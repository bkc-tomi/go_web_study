package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Emb は入力内容をuser.htmlに埋め込むためのもの
type Emb struct {
	Username string
	Age      int
	Gender   string
	Fruit    string
	Interest []string
}

// Msg は入力内容に誤りがある場合にlogin.htmlに埋め込むもの
type Msg struct {
	Message string
}

var emb Emb = Emb{}
var msg Msg = Msg{""}

func main() {
	port := "8080"
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/user", handleUser)
	log.Printf("server is running on http://localhost:%s/login", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 入力フォームを返す
		t, _ := template.ParseFiles("public/login.html")
		t.Execute(w, msg)
	}
	if r.Method == "POST" {
		msg.Message = ""
		r.ParseForm()
		// 必須入力のチェック
		if len(r.Form.Get("username")) == 0 {
			msg.Message = msg.Message + "usernameが入力されていません。"
		}
		// 値のチェック
		getint, interr := strconv.Atoi(r.Form.Get("age")) // 文字列の変換
		if interr != nil || getint < 0 || getint > 100 {
			msg.Message = msg.Message + "年齢がおかしいです。"
		}
		// ラジオボタンのチェック
		slice := []string{"1", "2"}
		exist := false
		for _, s := range slice {
			if r.Form.Get("gender") == s {
				exist = true
			}
		}
		if !exist {
			msg.Message = msg.Message + "性別がおかしいです。"
		}
		// プルダウンメニュー
		slice2 := []string{"apple", "pear", "banane"}
		exist2 := false
		for _, s := range slice2 {
			if r.Form.Get("fruit") == s {
				exist2 = true
			}
		}
		if !exist2 {
			msg.Message = msg.Message + "選択肢の中から選んでください。"
		}
		// チェックボックス
		slice3 := []string{"football", "baseball", "basketball", "tennis"}
		exist3 := 0
		interest := r.Form["interest"]
		for _, i := range interest {
			for _, s := range slice3 {
				if i == s {
					exist3++
				}
			}
		}
		if len(interest) != exist3 {
			msg.Message = msg.Message + "選択肢の中から選んでください。"
		}

		// 入力内容によってページ遷移
		if msg.Message == "" {
			emb = Emb{
				Username: r.Form.Get("username"),
				Age:      getint,
				Gender:   r.Form.Get("gender"),
				Fruit:    r.Form.Get("fruit"),
				Interest: r.Form["interest"],
			}
			http.Redirect(w, r, "/user", 301)
		} else {
			http.Redirect(w, r, "/login", 301)
		}
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("public/user.html")
	t.Execute(w, emb)
}
