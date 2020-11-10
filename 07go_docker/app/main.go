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
	Users   map[int]User
	Time    time.Time
}

// User db users
type User struct {
	ID       int
	Name     string
	Password string
}

const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	DataSourceName = "root:golang@tcp(mysql-container:3306)/golang_db"
)

var usr = make(map[int]User)
var templates = make(map[string]*template.Template)

func main() {
	// database
	db, dbErr := sql.Open(DriverName, DataSourceName)
	if dbErr != nil {
		log.Print("error connecting to database:", dbErr)
	}
	defer db.Close()
	rows, queryErr := db.Query("SELECT * FROM users")
	if queryErr != nil {
		log.Print("query error :", queryErr)
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Password); err != nil {
			log.Print(err)
		}
		usr[u.ID] = User{
			ID:       u.ID,
			Name:     u.Name,
			Password: u.Password,
		}
	}
	// web_server
	port := "8080"
	templates["index"] = loadTemplate("index")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handleIndex)
	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	temp := Embed{"Hello Golang!", "こんにちは！", usr, time.Now()}
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
