package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Emb struct {
	User string
	Art  string
}

var emb Emb

func main() {
	port := "8080"
	http.HandleFunc("/xss", handleXss)
	http.HandleFunc("/csrf", handleCsrf)
	http.HandleFunc("/view", handleView)
	log.Printf("server is running on http://localhost:%s/csrf", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}

func handleXss(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("public/xss.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Print(err)
			return
		}
		t.Execute(w, nil)
	}
	if r.Method == "POST" {
		r.ParseForm()
		user := r.Form.Get("username")
		if len(user) <= 0 {
			http.Redirect(w, r, "/xss", 301)
		} else {
			escape := template.HTMLEscapeString(user)
			emb = Emb{
				User: escape,
			}
			http.Redirect(w, r, "/view", 301)
		}
	}
}

func handleCsrf(w http.ResponseWriter, r *http.Request) {
	h := md5.New()
	io.WriteString(h, "golang")
	token := fmt.Sprintf("%x", h.Sum(nil))
	log.Print(token)
	if r.Method == "GET" {
		t, err := template.ParseFiles("public/csrf.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Print(err)
			return
		}
		t.Execute(w, token)
	}
	if r.Method == "POST" {
		r.ParseForm()
		requestToken := r.Form.Get("token")
		if requestToken != "" {
			// デベロッパーツールでtoken自体は見えてしまうのでもう一手間必要
			if token != requestToken {
				log.Print("tokenが一致しません。")
				return
			}
			user := r.Form.Get("username")
			art := r.Form.Get("article")
			if len(user) <= 0 && len(art) <= 0 {
				http.Redirect(w, r, "/csrf", 301)
			} else {
				emb = Emb{
					User: user,
					Art:  art,
				}
				log.Print("正しく処理されました。")
				http.Redirect(w, r, "/view", 301)
			}
		} else {
			log.Print("tokenがありません。")
		}
	}
}

func handleView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/view.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Print(err)
		return
	}
	t.Execute(w, emb)
}
