package main

import (
	"fmt"
	"net/http"
)

// 被害者の「ログイン済みセッション」を表す簡易Cookie。
const corsSessionCookie = "cors_session"

// Cookie(セッション)がある時だけ返す、victim本人の秘密データ。
const corsVictimSecret = `{"user":"alice","balance":"1234567","apiKey":"SECRET-9f3a2b"}`

// 許可リスト（対策版が許可する正規オリジン）。
const corsAllowedOrigin = "https://portal.corp.local"

// /cors : victim が普段使う正規ポータル画面。アクセスした時点で「ログイン済み」の
// Cookie を発行し（＝このブラウザは app の認証済みセッションを持つ状況を再現）、
// 本人には自分の秘密データAPIの場所を案内する。
func corsHandler(w http.ResponseWriter, r *http.Request) {
	// victim がこのポータルにログイン済み、という状況を再現（セッションCookie発行）
	http.SetCookie(w, &http.Cookie{Name: corsSessionCookie, Value: "victim-token-abc", Path: "/"})
	m := mode(r)
	body := fmt.Sprintf(`<h2>社内ポータル（victim視点）</h2>
<p>これは victim が普段使う正規サイト。<b>いま「ログイン済み」のCookieが発行された</b>（この後の攻撃で使われる）。</p>
<p>victim 本人の秘密データAPI: <code>GET /cors/api</code> … Cookieがある時だけ中身(残高やAPIキー)を返す。</p>
<p>現在のモード: <b>%s</b></p>
<hr>
<p style="color:#666">攻撃の本体は別オリジンの攻撃者ページ <code>http://127.0.0.1:9000/</code> で体験する。
この <code>/cors</code> を一度開いてログインCookieを作ってから、攻撃者ページを開くこと。</p>`, m)
	layout(w, "/cors", "CORSの設定ミス", m, body)
}

// /cors/api : victim の秘密データを返すAPI。セッションCookieがある時だけ中身を返す
// ＝「資格情報付き(credentialed)」リソース。別オリジンのJSがこの応答を読めるかは、
// 下のCORSヘッダの付け方次第で、最終的にブラウザが判定する。
func corsAPIHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	origin := r.Header.Get("Origin")

	// --- CORSヘッダの付け方（ここが脆弱/対策の差） ---
	if m == "vuln" {
		// 脆弱: 送られてきたOriginを無検証でそのまま反射し、credentials も許可する。
		// → どんな攻撃者オリジンでも「君に応答を読ませてよい」とブラウザに宣言してしまう。
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
	} else {
		// 対策: あらかじめ決めた正規オリジンに完全一致した時だけ許可する。
		if origin == corsAllowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
	}

	w.Header().Set("Content-Type", "application/json")

	// セッションCookieが無ければ秘密は返さない（credentialsが本質である点を示す）。
	if c, err := r.Cookie(corsSessionCookie); err != nil || c.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error":"not logged in (no session cookie)"}`)
		return
	}
	fmt.Fprint(w, corsVictimSecret)
}
