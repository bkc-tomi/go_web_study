package main

import (
	"app/mod/handler"
	_ "app/mod/memory"
	"app/mod/session"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var templates = make(map[string]*template.Template)

var globalSessions *session.Manager

//この後、init関数で初期化を行います。
func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func main() {
	// log.Print(handler.HelloHandler())
	// web_server
	port := "8080"
	// パース済みhtmlテンプレート
	templates["forum"] = handler.LoadTemplate("forum")
	templates["post"] = handler.LoadTemplate("post")
	templates["update"] = handler.LoadTemplate("update")

	// handlerの作成
	handleForum, fErr := handler.HandleForum(templates["forum"])
	if fErr != nil {
		log.Print(fErr)
	}
	handlePost, pErr := handler.HandlePost(templates["post"])
	if pErr != nil {
		log.Print(pErr)
	}
	handleDelete, dErr := handler.HandleDelete()
	if dErr != nil {
		log.Print(dErr)
	}
	handleEdit, eErr := handler.HandleEdit()
	if eErr != nil {
		log.Print(eErr)
	}
	handleUpdate, uErr := handler.HandleUpdate(templates["update"])
	if uErr != nil {
		log.Print(uErr)
	}
	// ルーティング
	http.HandleFunc("/", handleForum)
	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/edit", handleEdit)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/count", count)

	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}

var tokens []string

func count(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sess := globalSessions.SessionStart(w, r)
	t, _ := template.ParseFiles("public/count.html")
	w.Header().Set("Content-Type", "text/html")
	log.Println("sessionId: ", sess.SessionID())
	to, err := r.Cookie("token")
	if err != nil || to.Value == "" {
		// tokenが一致しなかった場合
		log.Println("no token.")
		globalSessions.SessionDestroy(w, r)
		sess = globalSessions.SessionStart(w, r)
		sess.Set("countnum", 1)
		token := createToken()
		tokens = append(tokens, token)
		cookie := http.Cookie{Name: "token", Value: token}
		http.SetCookie(w, &cookie)
		t.Execute(w, sess.Get("countnum"))
	} else {
		if contain(tokens, to.Value) {
			// tokenが一致した場合
			log.Println("contain tokens:", tokens, "get token:", to.Value)
			// ＊
			ct := sess.Get("countnum")
			if ct == nil {
				sess.Set("countnum", 1)
			} else {
				sess.Set("countnum", (ct.(int) + 1))
			}
			// 現在のtokenの削除
			tokens = removeToken(tokens, to.Value)
			token := createToken()
			tokens = append(tokens, token)
			cookie := http.Cookie{Name: "token", Value: token}
			http.SetCookie(w, &cookie)
			t.Execute(w, sess.Get("countnum"))
		} else {
			// tokenが一致しなかった場合
			log.Println("no contain tokens:", tokens, "get token:", to.Value)
			globalSessions.SessionDestroy(w, r)
			sess = globalSessions.SessionStart(w, r)
			sess.Set("countnum", 1)
			token := createToken()
			tokens = append(tokens, token)
			cookie := http.Cookie{Name: "token", Value: token}
			http.SetCookie(w, &cookie)
			t.Execute(w, sess.Get("countnum"))
		}
	}

	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 60) < (time.Now().Unix()) {
		globalSessions.SessionDestroy(w, r)
		sess = globalSessions.SessionStart(w, r)
	}

}

func contain(tokens []string, token string) bool {
	for _, to := range tokens {
		if to == token {
			return true
		}
	}
	return false
}

func createToken() string {
	h := md5.New()
	salt := "golang%^7&8888"
	io.WriteString(h, salt+time.Now().String())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func removeToken(tokens []string, token string) (newTokens []string) {
	for _, to := range tokens {
		if to != token {
			newTokens = append(newTokens, to)
		}
	}
	return
}
