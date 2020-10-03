package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"

	http.Handle("/", http.FileServer(http.Dir("root/")))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}
