package main

import (
	"net/http"
	"strings"
)

// オープンリダイレクト。遷移先を検証しないと任意サイトへ飛ばせる。
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	to := r.URL.Query().Get("to")
	if to != "" {
		if m == "vuln" {
			// 脆弱: 入力のURLへそのままリダイレクト
			http.Redirect(w, r, to, http.StatusFound)
			return
		}
		// 対策: 自サイト内（相対パス）のみ許可
		if strings.HasPrefix(to, "/") && !strings.HasPrefix(to, "//") {
			http.Redirect(w, r, to, http.StatusFound)
			return
		}
		http.Error(w, "許可されない遷移先です", http.StatusBadRequest)
		return
	}
	body := `<h2>リダイレクト</h2><p>?to=... で遷移。</p>
<hr><p style="color:#666">試:　<code>?to=https://example.com/</code>（脆弱版は外部へ飛ぶ／対策版は拒否）</p>`
	layout(w, "/redirect", "オープンリダイレクト", m, body)
}
