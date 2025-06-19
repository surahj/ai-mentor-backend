// test-server.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create a simple HTTP handler
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	// Start server
	fmt.Println("Simple test server starting on port 8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
