package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func main() {
	http.HandleFunc("/", markdownHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func markdownHandler(w http.ResponseWriter, r *http.Request) {
	out, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("iouutil.ReadAll: %w", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	unsafe := blackfriday.Run(out)
	// This is a very basic content policy and tighter standards are recommended.
	output := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	w.Write(output)
}
