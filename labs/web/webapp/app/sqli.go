package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
)

func sqliHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	q := r.URL.Query()
	u, p, search := q.Get("u"), q.Get("p"), q.Get("q")
	var b strings.Builder

	// --- ログイン（認証バイパスの題材） ---
	fmt.Fprintf(&b, `<h2>ログイン</h2><form><input type=hidden name=mode value=%s>
ユーザ名 <input name=u value="%s"> パスワード <input name=p value="%s"> <button>login</button></form>`,
		m, html.EscapeString(u), html.EscapeString(p))
	if u != "" {
		var name string
		var err error
		if m == "vuln" {
			// 脆弱: 文字列連結
			sqlText := fmt.Sprintf("SELECT username FROM users WHERE username='%s' AND password='%s'", u, p)
			err = db.QueryRow(sqlText).Scan(&name)
		} else {
			// 対策: プレースホルダ
			err = db.QueryRow("SELECT username FROM users WHERE username=? AND password=?", u, p).Scan(&name)
		}
		if err != nil {
			b.WriteString(`<p>ログイン失敗</p>`)
		} else {
			fmt.Fprintf(&b, `<p><b>ようこそ %s さん（成功）</b></p>`, html.EscapeString(name))
		}
	}

	// --- 商品検索（UNION抽出の題材） ---
	fmt.Fprintf(&b, `<h2>商品検索</h2><form><input type=hidden name=mode value=%s>
<input name=q value="%s"><button>search</button></form>`, m, html.EscapeString(search))
	if search != "" {
		var rows interface {
			Next() bool
			Scan(...any) error
			Close() error
		}
		var err error
		if m == "vuln" {
			rows, err = db.Query("SELECT name, price FROM products WHERE name LIKE '%" + search + "%'")
		} else {
			rows, err = db.Query("SELECT name, price FROM products WHERE name LIKE '%'||?||'%'", search)
		}
		if err != nil {
			fmt.Fprintf(&b, `<p>SQLエラー: %s</p>`, html.EscapeString(err.Error()))
		} else {
			b.WriteString("<table border=1><tr><th>name</th><th>price</th></tr>")
			for rows.Next() {
				var a, c string
				rows.Scan(&a, &c)
				fmt.Fprintf(&b, "<tr><td>%s</td><td>%s</td></tr>", html.EscapeString(a), html.EscapeString(c))
			}
			rows.Close()
			b.WriteString("</table>")
		}
	}

	b.WriteString(`<hr><p style="color:#666">試:　ユーザ名 <code>admin' -- </code>／検索 <code>x' UNION SELECT username,password FROM users -- </code></p>`)
	layout(w, "/sqli", "SQLインジェクション", m, b.String())
}
