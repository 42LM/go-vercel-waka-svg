// Package main implements a simple HTTP server with
// the single handler that is being used by vercel.
//
// The main package should only be used for local testing.
package main

import (
	"net/http"

	handler "go-vercel-waka-svg/api"
)

func main() {
	http.HandleFunc("/api", handler.GenerateSVG)

	http.ListenAndServe(":8080", nil)
}
