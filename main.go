package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.HandleFunc("/hi.txt", serveHiTxt)
	if lambdaEnv := os.Getenv("LAMBDA"); strings.ToLower(lambdaEnv) == "true" || os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Start Lambda handler
		Start(r)
	} else {
		// Start local HTTP server
		log.Println("Starting server on port 8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}
}

func serveHiTxt(w http.ResponseWriter, r *http.Request) {
	// Open the "./hi.txt" file
	data, err := os.ReadFile("./hi.txt")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}

	// Set the content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the file content to the response
	w.Write(data)
}
