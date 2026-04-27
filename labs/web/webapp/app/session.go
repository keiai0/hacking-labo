package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
)

var (
	sessMu   sync.Mutex
	sessions = map[string]string{} // sessionID -> username
)

func authUser(u, p string) bool {
	var x string
	err := db.QueryRow("SELECT username FROM users WHERE username=? AND password=?", u, p).Scan(&x)
	return err == nil
}

// セッション固定。ログイン時にセッションIDを再生成するかどうか。
func sessionHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	var b strings.Builder

	if r.Method == http.MethodPost {
		r.ParseForm()
		if authUser(r.FormValue("u"), r.FormValue("p")) {
			sid := ""
			if c, err := r.Cookie("SID"); err == nil {
				sid = c.Value
			}
			if m == "vuln" {
				// 脆弱: 提示されたセッションIDをそのまま認証済みにする（再生成しない）
				if sid == "" {
					sid = newToken()
				}
			} else {
				// 対策: ログイン時に必ず新しいIDを発行する
				sid = newToken()
			}
			sessMu.Lock()
			sessions[sid] = r.FormValue("u")
			sessMu.Unlock()
			http.SetCookie(w, &http.Cookie{Name: "SID", Value: sid, Path: "/"})
			fmt.Fprintf(&b, `<p>ログイン成功。セッションID = <code>%s</code></p>`, html.EscapeString(sid))
		} else {
			b.WriteString(`<p>ログイン失敗</p>`)
		}
	} else {
		// 現在のセッション状態
		if c, err := r.Cookie("SID"); err == nil {
			sessMu.Lock()
			u := sessions[c.Value]
			sessMu.Unlock()
			state := "未ログイン"
			if u != "" {
				state = "ログイン済み: " + u
			}
			fmt.Fprintf(&b, `<p>SID=<code>%s</code>（%s）</p>`, html.EscapeString(c.Value), html.EscapeString(state))
		}
	}

	fmt.Fprintf(&b, `<h2>ログイン</h2><form method=post action="/session?mode=%s">
ユーザ名 <input name=u> パスワード <input name=p> <button>login</button></form>`, m)
	b.WriteString(`<hr><p style="color:#666">セッション固定: 攻撃者が被害者に既知のSIDを事前に仕込む。
脆弱版はそのSIDがログイン後も有効になる／対策版は無効（新IDを発行）。</p>`)
	layout(w, "/session", "セッション固定", m, b.String())
}
