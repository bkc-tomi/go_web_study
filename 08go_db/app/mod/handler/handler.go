package handler

import (
	"app/mod/db"
	"html/template"
	"log"
	"net/http"
)

// HelloHandler モジュールの読み込み確認用
func HelloHandler() string {
	return "hello Handler"
}

// OriginalHandler パース済みのhtmlテンプレートと埋め込み構造体を受け取りハンドラー関数を返す
func OriginalHandler(t *template.Template, emb interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, emb); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
	}
}

func HandleForum(t *template.Template) (f func(w http.ResponseWriter, r *http.Request), err error) {
	data, err := db.GetTalks()

	if err != nil {
		return
	}
	f = func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, data); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
	}
	return
}

// LoadTemplate パース済みテンプレートを返す。
func LoadTemplate(name string) *template.Template {
	t, err := template.ParseFiles(
		"public/"+name+".html",
		"public/template/header.html",
		"public/template/footer.html",
	)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}
