// OAuth2/OIDC の設定・実装ミスを学ぶ軽量IdP+RP（Keycloak不使用、標準ライブラリのみ）。
// 各エンドポイントは ?mode=vuln|fixed で挙動を切り替える。
package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// 登録済みクライアント（RP）
const (
	clientID         = "corp-portal"
	registeredURI    = "http://client.corp/callback" // 正規のredirect_uri
	idTokenSecret    = "idp-signing-key"
	issuer           = "http://authz"
	audience         = "corp-portal"
)

type codeInfo struct {
	redirectURI string
	challenge   string // PKCE code_challenge（S256）
	sub         string
}

var (
	mu       sync.Mutex
	codes    = map[string]codeInfo{}
	states   = map[string]bool{} // RPが発行した有効なstate
	verifier = map[string]string{}
)

func mode(r *http.Request) string {
	if r.URL.Query().Get("mode") == "fixed" {
		return "fixed"
	}
	return "vuln"
}
func tok() string { b := make([]byte, 16); rand.Read(b); return hex.EncodeToString(b) }
func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
func unb64(s string) []byte { d, _ := base64.RawURLEncoding.DecodeString(s); return d }

func s256(v string) string { h := sha256.Sum256([]byte(v)); return b64(h[:]) }

func signJWT(alg, payload string) string {
	hdr := b64([]byte(fmt.Sprintf(`{"alg":"%s","typ":"JWT"}`, alg)))
	pl := b64([]byte(payload))
	if alg == "none" {
		return hdr + "." + pl + "."
	}
	mac := hmac.New(sha256.New, []byte(idTokenSecret))
	mac.Write([]byte(hdr + "." + pl))
	return hdr + "." + pl + "." + b64(mac.Sum(nil))
}

func main() {
	http.HandleFunc("GET /{$}", menu)
	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/token", tokenHandler)
	http.HandleFunc("/rp/login", rpLoginHandler)
	http.HandleFunc("/rp/callback", rpCallbackHandler)
	http.HandleFunc("/rp/validate", validateHandler)
	http.HandleFunc("/attacker/catch", attackerHandler)
	log.Println("listening :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func menu(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<!doctype html><meta charset="utf-8"><title>OAuth/OIDC lab</title>
<body style="font-family:sans-serif;max-width:820px;margin:2rem auto">
<h1>OAuth2 / OIDC 設定ミス演習（軽量IdP+RP）</h1>
<p>登録済み: client_id=<code>corp-portal</code> / redirect_uri=<code>http://client.corp/callback</code></p>
<ol>
<li><b>redirect_uri 検証不備</b>: <code>/authorize?mode=vuln&client_id=corp-portal&redirect_uri=...&response_type=code</code></li>
<li><b>state 欠如(CSRF)</b>: <code>/rp/login?mode=vuln</code> と <code>/rp/callback</code></li>
<li><b>PKCE なし</b>: <code>/authorize</code>+<code>/token?mode=vuln</code></li>
<li><b>IDトークン検証不備</b>: <code>/rp/validate?mode=vuln&token=...</code></li>
</ol>
<p style="color:#666">詳細な手順と curl 例は README を参照。各エンドポイントは ?mode=vuln|fixed。</p>
</body>`)
}
