// Package design goa service DSL
package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("private", func() {
	Description("This service provides private functions.")

	Method("ping", func() {
		Result(String)

		HTTP(func() {
			GET("/ping")
		})

		GRPC(func() {
		})
	})

})
