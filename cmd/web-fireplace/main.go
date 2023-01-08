package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	rootDir = "."
)

func main() {
	fmt.Println("Web Fireplace")

	fireSrv := http.FileServer(http.Dir(rootDir))

	if err := http.ListenAndServe(":8080", fireSrv); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
