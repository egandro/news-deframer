package deframer

import (
	"testing"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/openai"
	"github.com/egandro/news-deframer/pkg/openai/mock_openai"
	"github.com/egandro/news-deframer/pkg/source"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTestDeframer(t *testing.T, ai openai.OpenAI, src *source.Source) (Deframer, error) {
	// Use in-memory SQLite for testing
	db, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	res := &deframer{
		db:  db,
		ai:  ai,
		src: src,
	}

	return res, nil
}

func TestNewDeframer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	openAIMock := mock_openai.NewMockOpenAI(ctrl)

	s, err := setupTestDeframer(t, openAIMock, nil)
	//s, err := NewDeframer()
	assert.NoError(t, err)
	assert.NotNil(t, s, "Deframer should be initialized")
}
