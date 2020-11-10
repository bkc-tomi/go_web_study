package main

import (
	"fmt"
	"log"
	"net/http"
)

// MyMux is Original ServeMux
type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		handleIndex(w, r)
	case "/greeting":
		handleGreet(w, r)
	default:
		http.NotFound(w, r)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome golang server.")
}
func handleGreet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello guest!")
}

func main() {
	port := "8080"

	mux := &MyMux{}
	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, mux))
}
