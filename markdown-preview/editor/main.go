package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(0)

	var err error
	s, err := NewServiceFromEnv()
	if err != nil {
		log.Fatalf("NewServiceFromEnv %v", err)
	}
	mux := s.RegisterHandlers()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
