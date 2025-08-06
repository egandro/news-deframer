// Package design goa service DSL
package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("web", func() {
	Description("Web service that returns HTML content")

	Error("invalid_feed_id", String, "Invalid Feed Id")

	HTTP(func() {
		Response("invalid_feed_id", StatusNotFound)
	})

	Method("index", func() {
		Description("Returns the index page in HTML")

		HTTP(func() {
			GET("/")
			Response(StatusOK, func() {
				ContentType("text/html")
			})
		})

		Result(String)
	})

	Method("feed", func() {
		Description("Returns the feed with the given xml")

		Payload(func() {
			Attribute("feed_id", String, "Feed Id", func() {
				Example("some-id")
			})
			Required("feed_id")
		})

		HTTP(func() {
			GET("/feed/{feed_id}")
			SkipResponseBodyEncodeDecode()
			Response(func() {
				Header("length:Content-Length") // Map length to Content-Length header
				Header("type:Content-Type")     // Map type to Content-Type header
			})
		})

		Error("invalid_feed_id")

		Result(func() {
			// We'll return the file size in the Content-Length header
			Attribute("length", Int64, "Content length in bytes")
			Attribute("type", String, "Content type")
			Required("length", "type")
		})

	})
})
