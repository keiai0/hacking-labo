package main

import (
	"net/http"
)

// CORSの設定ミス。Originを無検証で反射し、かつ資格情報を許可すると危険。
func corsHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	origin := r.Header.Get("Origin")
	if m == "vuln" {
		// 脆弱: 任意のOriginを反射し、credentialsも許可
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
	} else {
		// 対策: allowlistに一致したときだけ許可
		if origin == "https://portal.corp.local" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
	}
	body := `<h2>CORS API</h2><p>このページのレスポンスヘッダを確認。</p>
<hr><p style="color:#666">試:　<code>curl -H 'Origin: https://evil.example' -i …/cors</code>
　脆弱版はそのOriginを反射する。</p>`
	layout(w, "/cors", "CORSの設定ミス", m, body)
}
