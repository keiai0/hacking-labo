package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html"
	"net/http"
	"strings"
)

var currentEmail = "alice@corp.local"
var csrfToken = newToken()

func newToken() string {
	bts := make([]byte, 16)
	rand.Read(bts)
	return hex.EncodeToString(bts)
}

// CSRF。状態変更に正規トークンを要求するかどうか。
func csrfHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	var b strings.Builder

	if r.Method == http.MethodPost {
		r.ParseForm()
		email := r.FormValue("email")
		if m == "fixed" && r.FormValue("token") != csrfToken {
			// 対策: トークン不一致は拒否
			http.Error(w, "CSRFトークン不正", http.StatusForbidden)
			return
		}
		currentEmail = email
		fmt.Fprintf(&b, `<p>メールアドレスを <b>%s</b> に変更しました</p>`, html.EscapeString(email))
	}

	tokenField := ""
	if m == "fixed" {
		tokenField = fmt.Sprintf(`<input type=hidden name=token value="%s">`, csrfToken)
	}
	fmt.Fprintf(&b, `<h2>メール変更</h2><p>現在: %s</p>
<form method=post action="/csrf?mode=%s">%s
新アドレス <input name=email><button>変更</button></form>`,
		html.EscapeString(currentEmail), m, tokenField)

	b.WriteString(`<hr><p style="color:#666">攻撃イメージ: 外部サイトから<code>POST /csrf</code>を送る。
脆弱版はトークンなしでも通る／対策版は拒否(403)。</p>`)
	layout(w, "/csrf", "CSRF", m, b.String())
}
