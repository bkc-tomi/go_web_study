package main

import (
	"log"
	"net/http"
)

func main() {
	port := "4000"
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	log.Printf("Server listening on http://localhost:%s/", port)
	http.ListenAndServe(":"+port, nil)
}
