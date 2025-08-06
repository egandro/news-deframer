package main

import (
	"context"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/deframer"
	"github.com/joho/godotenv"
	"goa.design/clue/log"
)

const deleteOrphanedJobsTime = 10 * time.Minute

// bootstrap our own services
func bootstrap(ctx context.Context, httpPortF *string, dbgF *bool) (outHttpPortF *string, outDbgF *bool) {
	outHttpPortF = httpPortF
	outDbgF = dbgF

	_ = godotenv.Load() // load .env file - if exist
	cfg, err := config.GetConfig()

	if err != nil {
		log.Fatalf(ctx, err, "can't initialize config")
	}

	if cfg.HttpPort != "" {
		outHttpPortF = &cfg.HttpPort
	}

	if cfg.DebugLog {
		*outDbgF = true
	}

	d, err := deframer.NewDeframer(ctx)
	if err != nil {
		log.Fatalf(ctx, err, "can't create deframer")
	}

	_, err = d.UpdateFeeds()
	if err != nil {
		log.Fatalf(ctx, err, "can't update feeds")
	}

	return
}
