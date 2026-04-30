package main

import (
	"fmt"
	"net/http"
	"strings"
)

// 1) 認可エンドポイント: redirect_uri の検証が主題
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	m := mode(r)
	cid := q.Get("client_id")
	ruri := q.Get("redirect_uri")
	state := q.Get("state")
	challenge := q.Get("code_challenge")

	if cid != clientID {
		http.Error(w, "unknown client", http.StatusBadRequest)
		return
	}

	ok := false
	if m == "vuln" {
		// 脆弱: 前方一致で検証 → client.corp.evil.example 等が通る
		ok = strings.HasPrefix(ruri, "http://client.corp")
	} else {
		// 対策: 登録値と完全一致
		ok = ruri == registeredURI
	}
	if !ok {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}

	code := tok()
	mu.Lock()
	codes[code] = codeInfo{redirectURI: ruri, challenge: challenge, sub: "alice"}
	mu.Unlock()

	loc := ruri + "?code=" + code
	if state != "" {
		loc += "&state=" + state
	}
	http.Redirect(w, r, loc, http.StatusFound) // Location に認可コードが乗る
}

// 3) トークンエンドポイント: PKCE検証が主題
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	r.ParseForm()
	get := func(k string) string {
		if v := r.FormValue(k); v != "" {
			return v
		}
		return r.URL.Query().Get(k)
	}
	code := get("code")
	ruri := get("redirect_uri")
	cv := get("code_verifier")

	mu.Lock()
	info, found := codes[code]
	mu.Unlock()
	if !found {
		http.Error(w, "invalid code", http.StatusBadRequest)
		return
	}
	if ruri != info.redirectURI {
		http.Error(w, "redirect_uri mismatch", http.StatusBadRequest)
		return
	}
	// PKCE
	if info.challenge != "" {
		if m == "fixed" {
			if cv == "" || s256(cv) != info.challenge {
				http.Error(w, "PKCE検証失敗", http.StatusBadRequest)
				return
			}
		}
		// vuln: code_verifier を検証しない（横取りコードがそのまま使える）
	}

	mu.Lock()
	delete(codes, code)
	mu.Unlock()

	idToken := signJWT("HS256", fmt.Sprintf(`{"iss":"%s","aud":"%s","sub":"%s","role":"user"}`, issuer, audience, info.sub))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"access_token":"%s","token_type":"Bearer","id_token":"%s"}`, tok(), idToken)
}

// 攻撃者の受け皿（横取りした認可コードが届く様子）
func attackerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "攻撃者サーバが認可コードを受信: %s\n", r.URL.Query().Get("code"))
}
