package handler

import (
	"app/mod/db"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
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

var datas struct {
	Talks []db.Talk
}

// HandleForum パース済みのhtmlテンプレートと埋め込み構造体を受け取り,掲示板ページへのハンドラーを返す
func HandleForum(t *template.Template) (f func(w http.ResponseWriter, r *http.Request), err error) {
	f = func(w http.ResponseWriter, r *http.Request) {
		// dbからデータの取得
		talks, err := db.GetTalks()
		datas.Talks = talks
		if err != nil {
			return
		}
		// cookieの設定
		// expiration := time.Now()
		// expiration = expiration.AddDate(0, 0, 1)
		// cookie := http.Cookie{Name: "username", Value: "golang", Expires: expiration}
		// http.SetCookie(w, &cookie)

		// for _, c := range r.Cookies() {
		// 	log.Print("Name:", c.Name, ", Value:", c.Value)
		// }

		if err := t.Execute(w, datas); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
	}
	return
}

var msg struct {
	Message string
}

// HandlePost パース済みのhtmlテンプレートと埋め込み構造体を受け取り,投稿ページへのハンドラーを返す
func HandlePost(t *template.Template) (f func(w http.ResponseWriter, r *http.Request), err error) {
	f = func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			// log.Print("/post GET")
			t.Execute(w, msg)
			msg.Message = ""
		}
		if r.Method == "POST" {
			// log.Print("/post POST")
			r.ParseForm()
			l := r.Form.Get("posttext")
			if len(l) <= 0 || len(l) > 140 {
				// 投稿内容が正しくない場合
				msg.Message = "文字数は１〜１４０字の間でお願いします。"
				http.Redirect(w, r, "/post", 301)
				// log.Print("talk length not match redirect /post")
			} else {
				// 投稿内容が正しい場合
				if err := db.PostTalk(l); err != nil {
					msg.Message = err.Error()
					http.Redirect(w, r, "/post", 301)
					// log.Print("db write err redirect /post")
				} else {
					http.Redirect(w, r, "/", 301)
					// log.Print("success redirect /")
				}
			}
		}
	}
	return
}

// HandleDelete 指定したIDのトークを削除する。
func HandleDelete() (f func(w http.ResponseWriter, r *http.Request), err error) {
	f = func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			id, err := strconv.Atoi(r.Form.Get("delete_id"))
			if err != nil {
				http.Redirect(w, r, "/", 301)
			}
			err = db.DeleteTalk(id)
			if err != nil {
				http.Redirect(w, r, "/", 301)
			}
			http.Redirect(w, r, "/", 301)
		}
	}
	return
}

// T 編集の際のトークの内容を保持
var T db.Talk

// HandleEdit 編集用のハンドラーを返す
func HandleEdit() (f func(w http.ResponseWriter, r *http.Request), err error) {
	f = func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			id, err := strconv.Atoi(r.Form.Get("edit_id"))
			if err != nil {
				http.Redirect(w, r, "/", 301)
			}
			t, err := db.GetTalk(id)
			T = t
			http.Redirect(w, r, "/update", 301)
		}
	}
	return
}

// HandleUpdate 指定したIDのトークの変更を上書きする。
func HandleUpdate(t *template.Template) (f func(w http.ResponseWriter, r *http.Request), err error) {
	f = func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			t.Execute(w, T)
			T = db.Talk{}
		}
		if r.Method == "POST" {
			r.ParseForm()

			id, _ := strconv.Atoi(r.Form.Get("update_id"))
			t := r.Form.Get("talk")

			talk, _ := db.GetTalk(id)
			talk.Talk = t
			talk.UpdateAt = time.Now()
			err := db.UpdateTalk(talk)
			if err != nil {
				log.Print(err)
				http.Redirect(w, r, "/", 301)
			}
			http.Redirect(w, r, "/", 301)
		}
	}
	return
}

// LoadTemplate パース済みテンプレートを返す。
func LoadTemplate(name string) *template.Template {
	t, err := template.ParseFiles(
		"public/" + name + ".html",
	)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}
