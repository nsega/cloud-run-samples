package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Verify the dot utility is available at startup
	// instead of waiting for a first request.
	fileInfo, err := os.Stat("/usr/bin/dot")
	if err != nil {
		log.Fatalf("graphviz-web: %v", err)
	}
	if fileInfo.Mode()&0111 == 0 {
		log.Fatal(`graphviz-web: ("/usr/bin/dot") not executable`)
	}

	http.HandleFunc("/diagram.png", diagramHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// diagramHandler renders a diagram using HTTP request parameters and the dot command.
func diagramHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("method not allowed: %s", r.Method)
		http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	dot := q.Get("dot")
	if dot == "" {
		log.Print("no graphviz definition provided")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Cache header must be set before writing a response.
	w.Header().Set("Cache-Control", "public, max-age=86400")

	input := strings.NewReader(dot)
	if err := createDiagram(w, input); err != nil {
		log.Printf("createDiagram: %v", err)
		// Do not cache error responses.
		w.Header().Del("Cache-Control")
		if strings.Contains(err.Error(), "syntax") {
			http.Error(w, "Bad Request: DOT syntax error", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// createDiagram generates a diagram image from the provided io.Reader written to the io.Writer.
func createDiagram(w io.Writer, r io.Reader) error {
	stderr := new(bytes.Buffer)
	args := []string{
		"-Glabel=Made on Cloud Run",
		"-Gfontsize=10",
		"-Glabeljust=right",
		"-Glabelloc=bottom",
		"-Gfontcolor=gray",
		"-Tpng",
	}
	cmd := exec.Command("/usr/bin/dot", args...)
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec(%s) failed (%v): %s", cmd.Path, err, stderr.String())
	}
	return nil
}
