package main

import (
	"github.com/richjyoung/echopen"
)

const Description = `Minimal example to get the server running`

func main() {
	// Create a new echOpen wrapper
	api := echopen.New("Minimal", "1.0.0", "3.1.0")
	api.Description(Description)

	// Serve the generated schema
	api.ServeSchema("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Start the server
	api.Start("localhost:3030")
}
