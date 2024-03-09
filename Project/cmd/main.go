package main

import (
	"log"
	"net/http"

	"github.com/E4kere/simpleproject/pkg/modules"
)

func main() {
	// Initialize the CS2 stats module
	modules.Init()

	// Register HTTP handlers
	http.HandleFunc("/players", modules.PlayerHandler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
