package main

import (
	"log"
	"net/http"
)

func main() {
	port := "3000"
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}
