package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
)

// 反射型XSS。入力をそのままHTMLに出すか、エスケープするかの違い。
func xssHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	name := r.URL.Query().Get("name")
	var b strings.Builder

	fmt.Fprintf(&b, `<h2>あいさつ</h2><form><input type=hidden name=mode value=%s>
名前 <input name=name value="%s"><button>送信</button></form>`, m, html.EscapeString(name))

	if name != "" {
		if m == "vuln" {
			// 脆弱: 入力をそのままHTMLに連結（html/templateを通さない）→ スクリプトが動く
			fmt.Fprintf(&b, `<p>こんにちは、%s さん</p>`, name)
		} else {
			// 対策: HTMLエスケープしてから出力
			fmt.Fprintf(&b, `<p>こんにちは、%s さん</p>`, html.EscapeString(name))
		}
	}

	b.WriteString(`<hr><p style="color:#666">試:　名前欄に <code>&lt;script&gt;alert(1)&lt;/script&gt;</code> や <code>&lt;img src=x onerror=alert(document.domain)&gt;</code></p>`)
	layout(w, "/xss", "反射型XSS", m, b.String())
}
