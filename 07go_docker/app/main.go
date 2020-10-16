package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Embed htmlファイルに埋め込むデータ構造体
type Embed struct {
	Title   string
	Message string
	Time    time.Time
}

const (
	DRIVER_NAME = "mysql" // ドライバ名(mysql固定)
	// user:password@tcp(container-name:port)/dbname ※mysql はデフォルトで用意されているDB
	DATA_SOURCE_NAME = "root:golang@tcp(mysql-container:3306)/mysql"
)

var templates = make(map[string]*template.Template)

func main() {
	// database
	db, err := sql.Open(DRIVER_NAME, DATA_SOURCE_NAME)
	if err != nil {
		log.Print("error connecting to database:", err)
	}
	log.Print(db)
	// web_server
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
