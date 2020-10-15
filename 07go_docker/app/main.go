package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// Embed htmlファイルに埋め込むデータ構造体
type Embed struct {
	Title   string
	Message string
	Time    time.Time
}

var templates = make(map[string]*template.Template)

func main() {
	port := "8080"

	templates["index"] = loadTemplate("index")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/index", handleIndex)
	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	temp := Embed{"Hello Golang!", "こんにちは！", time.Now()}
	if err := templates["index"].Execute(w, temp); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles(
		"root/"+name+".html",
		"root/template/header.html",
		"root/template/footer.html",
	)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}
