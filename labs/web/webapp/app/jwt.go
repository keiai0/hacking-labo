package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
)

var jwtSecret = []byte("supersecret-signing-key")

func b64(b []byte) string      { return base64.RawURLEncoding.EncodeToString(b) }
func unb64(s string) []byte    { d, _ := base64.RawURLEncoding.DecodeString(s); return d }
func sign(signingInput string) string {
	mac := hmac.New(sha256.New, jwtSecret)
	mac.Write([]byte(signingInput))
	return b64(mac.Sum(nil))
}

func makeJWT(role string) string {
	h := b64([]byte(`{"alg":"HS256","typ":"JWT"}`))
	p := b64([]byte(fmt.Sprintf(`{"user":"guest","role":"%s"}`, role)))
	return h + "." + p + "." + sign(h+"."+p)
}

// JWTの検証不備。alg=none や署名未検証を受け入れるか。
func jwtHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	token := r.URL.Query().Get("token")
	var b strings.Builder
	fmt.Fprintf(&b, `<p>配布トークン(role=guest): <code>%s</code></p>`, makeJWT("guest"))
	fmt.Fprintf(&b, `<form><input type=hidden name=mode value=%s>
token <input name=token value="%s" size=60><button>検証</button></form>`, m, html.EscapeString(token))

	if token != "" {
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			b.WriteString(`<p>形式不正</p>`)
		} else {
			var hdr struct{ Alg string }
			var pl struct{ User, Role string }
			json.Unmarshal(unb64(parts[0]), &hdr)
			json.Unmarshal(unb64(parts[1]), &pl)
			valid := false
			switch {
			case m == "vuln" && hdr.Alg == "none":
				valid = true // 脆弱: alg=none を無検証で受理
			case hdr.Alg == "HS256" && hmac.Equal([]byte(parts[2]), []byte(sign(parts[0]+"."+parts[1]))):
				valid = true // 署名が正しいHS256
			}
			if !valid {
				b.WriteString(`<p>検証失敗（拒否）</p>`)
			} else if pl.Role == "admin" {
				b.WriteString(`<p style="color:#b00"><b>管理者としてアクセス許可</b></p>`)
			} else {
				fmt.Fprintf(&b, `<p>role=%s として許可</p>`, html.EscapeString(pl.Role))
			}
		}
	}
	b.WriteString(`<hr><p style="color:#666">攻撃: ヘッダを<code>{"alg":"none"}</code>、ペイロードを<code>{"role":"admin"}</code>にし署名を空にする。
脆弱版はこれを管理者として受理／対策版は拒否。</p>`)
	layout(w, "/jwt", "JWTの検証不備", m, b.String())
}
