package main

import (
	"app/mod/handler"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Username string
}

var templates = make(map[string]*template.Template)

func main() {
	log.Print(handler.HelloHandler())
	// web_server
	port := "8080"
	// パース済みhtmlテンプレート
	templates["forum"] = handler.LoadTemplate("forum")
	templates["edit"] = handler.LoadTemplate("edit")
	templates["post"] = handler.LoadTemplate("post")

	// handlerの作成
	handleForum, fErr := handler.HandleForum(templates["forum"])
	if fErr != nil {
		log.Print(fErr)
	}
	handleEdit := handler.OriginalHandler(templates["edit"], nil)
	handlePost := handler.OriginalHandler(templates["post"], nil)
	// ルーティング
	http.HandleFunc("/", handleForum)
	http.HandleFunc("/edit", handleEdit)
	http.HandleFunc("/post", handlePost)
	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}
