package main

import (
	"app/restapi"
	"app/statichost"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	port := "8080"
	apiVersion := "v1"
	router := mux.NewRouter().StrictSlash(true)
	// static
	router.HandleFunc("/", statichost.HandleIndex)
	router.HandleFunc("/create", statichost.HandleCreate)
	router.HandleFunc("/edit", statichost.HandleEdit)
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("static/js/"))))
	// api
	router.HandleFunc("/api/"+apiVersion+"/users", restapi.GetUsers).Methods("GET")
	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, router))
}
