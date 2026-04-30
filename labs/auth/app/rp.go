package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

// 2) RP: ログイン開始。state / PKCE の有無が主題
func rpLoginHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	authz := "/authorize?mode=" + m + "&client_id=" + clientID +
		"&redirect_uri=" + registeredURI + "&response_type=code"
	if m == "fixed" {
		st := tok()
		mu.Lock()
		states[st] = true
		mu.Unlock()
		http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: st, Path: "/"})
		authz += "&state=" + st
		// PKCE（デモのため固定のverifier/challenge）
		authz += "&code_challenge=" + s256("verifier-xyz") + "&code_challenge_method=S256"
	}
	// vuln: state も PKCE も付けない
	fmt.Fprintf(w, "RPは次のURLへリダイレクトします:\n%s\n", authz)
}

// 2) RP: コールバック。state 検証が主題
func rpCallbackHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if m == "fixed" {
		c, err := r.Cookie("oauth_state")
		mu.Lock()
		valid := err == nil && c.Value == state && states[state]
		mu.Unlock()
		if !valid {
			http.Error(w, "state検証失敗（CSRFの疑い）", http.StatusBadRequest)
			return
		}
	}
	// vuln: state を検証しない → 攻撃者が用意したcodeで被害者をログインさせられる
	fmt.Fprintf(w, "ログイン処理: code=%s を受理（mode=%s）\n", code, m)
}

// 4) IDトークン検証エンドポイント
func validateHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	token := r.URL.Query().Get("token")
	parts := splitJWT(token)
	if parts == nil {
		http.Error(w, "形式不正", http.StatusBadRequest)
		return
	}
	var hdr struct{ Alg string }
	var pl struct{ Iss, Aud, Sub, Role string }
	json.Unmarshal(unb64(parts[0]), &hdr)
	json.Unmarshal(unb64(parts[1]), &pl)

	valid := false
	if m == "vuln" {
		// 脆弱: 署名も iss/aud も検証しない（alg=none も受理）
		valid = true
	} else {
		// 対策: 署名(HS256) と iss/aud を検証
		mac := hmac.New(sha256.New, []byte(idTokenSecret))
		mac.Write([]byte(parts[0] + "." + parts[1]))
		sigOK := hdr.Alg == "HS256" && hmac.Equal([]byte(parts[2]), []byte(b64(mac.Sum(nil))))
		valid = sigOK && pl.Iss == issuer && pl.Aud == audience
	}
	if !valid {
		http.Error(w, "IDトークン検証失敗（拒否）", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "受理: sub=%s role=%s\n", pl.Sub, pl.Role)
}

func splitJWT(t string) []string {
	n, parts := 0, make([]string, 0, 3)
	start := 0
	for i := 0; i <= len(t); i++ {
		if i == len(t) || t[i] == '.' {
			parts = append(parts, t[start:i])
			start = i + 1
			n++
		}
	}
	if n != 3 {
		return nil
	}
	return parts
}
