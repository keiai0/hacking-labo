package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SSRF。サーバがユーザ指定URLを取得すると、内部資源に到達できる。
func ssrfHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	target := r.URL.Query().Get("url")
	var b strings.Builder
	fmt.Fprintf(&b, `<h2>URLプレビュー</h2><form><input type=hidden name=mode value=%s>
URL <input name=url value="%s" size=40><button>取得</button></form>`, m, html.EscapeString(target))

	if target != "" {
		allowed := true
		if m == "fixed" {
			// 対策: スキームとホストをallowlistで検証（内部や非httpを拒否）
			u, err := url.Parse(target)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() != "example.com" {
				allowed = false
			}
		}
		if !allowed {
			b.WriteString(`<p>許可されていないURLです</p>`)
		} else {
			// 脆弱版: 任意URLをサーバが取得（内部10.10.5.x やメタデータにも到達しうる）
			c := &http.Client{Timeout: 3 * time.Second}
			resp, err := c.Get(target)
			if err != nil {
				fmt.Fprintf(&b, `<p>取得失敗: %s</p>`, html.EscapeString(err.Error()))
			} else {
				body, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
				resp.Body.Close()
				fmt.Fprintf(&b, `<p>HTTP %d</p><pre>%s</pre>`, resp.StatusCode, html.EscapeString(string(body)))
			}
		}
	}
	b.WriteString(`<hr><p style="color:#666">試:　<code>http://app:8080/sqli</code> など内部URL（脆弱版はサーバ側から到達）</p>`)
	layout(w, "/ssrf", "SSRF", m, b.String())
}
