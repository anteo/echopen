package main

import (
	"github.com/richjyoung/echopen"
)

func main() {
	// Create a new echOpen wrapper
	api := echopen.New(
		"Minimal",
		"1.0.0",
		echopen.WithSpecDescription("Minimal example to get the server running."),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeSwaggerUI("/", "/openapi.yml", "5.10.3")

	// Write the full generated spec
	api.WriteYAMLSpec("openapi_out.yml")

	// Start the server
	api.Start("localhost:3000")
}
