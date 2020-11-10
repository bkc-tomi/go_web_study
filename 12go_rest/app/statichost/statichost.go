package statichost

import (
	"html/template"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, nil)
}

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/create.html")
	t.Execute(w, nil)
}

func HandleEdit(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/edit.html")
	t.Execute(w, nil)
}
