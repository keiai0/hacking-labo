// 統合脆弱Webアプリ。1ページ=1脆弱性、各ページ ?mode=vuln|fixed で切替。
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// メニュー（ページを追加したらここに1行足す）
var menu = []struct{ Path, Name string }{
	{"/sqli", "SQLインジェクション"},
	{"/xss", "反射型XSS"},
	{"/cmdi", "OSコマンドインジェクション"},
	{"/path", "パストラバーサル"},
	{"/ssrf", "SSRF"},
	{"/redirect", "オープンリダイレクト"},
	{"/cors", "CORS設定ミス"},
	{"/csrf", "CSRF"},
}

// mode は vuln（既定）か fixed
func mode(r *http.Request) string {
	if r.URL.Query().Get("mode") == "fixed" {
		return "fixed"
	}
	return "vuln"
}

// ページ共通の枠（メニュー＋モード切替＋本文）
func layout(w http.ResponseWriter, path, title, m, body string) {
	fmt.Fprintf(w, `<!doctype html><meta charset="utf-8"><title>%s</title>
<body style="font-family:sans-serif;max-width:760px;margin:2rem auto">
<nav style="font-size:.9em">`, title)
	for _, it := range menu {
		fmt.Fprintf(w, `<a href="%s">%s</a>&nbsp;|&nbsp;`, it.Path, it.Name)
	}
	fmt.Fprintf(w, `</nav><hr><h1>%s</h1>`, title)
	if path != "" {
		other := "vuln"
		if m == "vuln" {
			other = "fixed"
		}
		cur := "脆弱版"
		if m == "fixed" {
			cur = "対策版"
		}
		fmt.Fprintf(w, `<p>現在: <b>%s</b>（<a href="%s?mode=%s">%s版に切替</a>）</p>`,
			cur, path, other, map[string]string{"vuln": "脆弱", "fixed": "対策"}[other])
	}
	fmt.Fprint(w, body, "</body>")
}

func main() {
	var err error
	db, err = sql.Open("sqlite", "/tmp/app.db")
	if err != nil {
		log.Fatal(err)
	}
	seed()

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		layout(w, "", "社内ポータル（演習用）", "", "<p>上部メニューから演習ページを選択。</p>")
	})
	http.HandleFunc("/sqli", sqliHandler)
	http.HandleFunc("/xss", xssHandler)
	http.HandleFunc("/cmdi", cmdiHandler)
	http.HandleFunc("/path", pathHandler)
	http.HandleFunc("/ssrf", ssrfHandler)
	http.HandleFunc("/redirect", redirectHandler)
	http.HandleFunc("/cors", corsHandler)
	http.HandleFunc("/csrf", csrfHandler)

	log.Println("listening :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func seed() {
	db.Exec(`CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY, username TEXT, password TEXT)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS products(id INTEGER PRIMARY KEY, name TEXT, price TEXT)`)
	db.Exec(`DELETE FROM users`)
	db.Exec(`DELETE FROM products`)
	db.Exec(`INSERT INTO users(username,password) VALUES('admin','S3cr3t-admin!'),('alice','alicepw')`)
	db.Exec(`INSERT INTO products(name,price) VALUES('Keyboard','3000'),('Mouse','1500'),('Monitor','20000')`)
}
