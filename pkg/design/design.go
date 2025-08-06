// Package design goa service DSL
package design

import (
	. "goa.design/goa/v3/dsl"

	// added for cors
	cors "goa.design/plugins/v3/cors/dsl"
)

// this is a regex
const corsHeader = "/.*localhost.*/"

var _ = API("service", func() {
	Title("Service")

	Description("Deframe RSS Feeds")

	Server("service", func() {
		Host("default", func() {
			Description("Default hosts.")
			// Transport specific URLs, supported schemes are:
			// 'http', 'https', 'grpc' and 'grpcs' with the respective default
			// ports: 80, 443, 8080, 8443.
			URI("http://0.0.0.0:8000")
			URI("grpc://0.0.0.0:8080")
		})
	})

	// added for cors (swagger needs this)
	cors.Origin(corsHeader, func() {
		cors.Headers("Content-Type", "api_key", "Authorization")
		cors.Methods("GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS")
		cors.MaxAge(600)
	})
})
