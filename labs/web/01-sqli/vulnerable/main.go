// 脆弱版: SQL文を文字列連結で組み立てている（SQLインジェクション）
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

var db *sql.DB

var page = template.Must(template.New("p").Parse(`<!doctype html><meta charset="utf-8">
<title>Shop (vulnerable)</title>
<h1>社員ポータル（脆弱版）</h1>
<h2>ログイン</h2>
<form method="post" action="/login">
  ユーザ名 <input name="username">
  パスワード <input name="password" type="password">
  <button>ログイン</button>
</form>
<p>{{.Msg}}</p>
<h2>商品検索</h2>
<form method="get" action="/search">
  <input name="q" value="{{.Q}}"><button>検索</button>
</form>
{{if .Rows}}<table border=1><tr><th>name</th><th>price</th></tr>
{{range .Rows}}<tr><td>{{index . 0}}</td><td>{{index . 1}}</td></tr>{{end}}</table>{{end}}
`))

type view struct {
	Msg  string
	Q    string
	Rows [][2]string
}

func render(w http.ResponseWriter, v view) { page.Execute(w, v) }

func main() {
	var err error
	db, err = sql.Open("sqlite", "/tmp/app.db")
	if err != nil {
		log.Fatal(err)
	}
	seed()

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		render(w, view{})
	})

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		u, p := r.FormValue("username"), r.FormValue("password")
		// 脆弱: 入力を文字列連結でSQLに埋め込む
		q := fmt.Sprintf("SELECT username FROM users WHERE username='%s' AND password='%s'", u, p)
		log.Println("login query:", q)
		var name string
		err := db.QueryRow(q).Scan(&name)
		if err != nil {
			render(w, view{Msg: "ログイン失敗"})
			return
		}
		render(w, view{Msg: "ようこそ、" + name + " さん（ログイン成功）"})
	})

	http.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		// 脆弱: LIKE句に入力を直接連結（UNIONベースの抽出が可能）
		sqlText := "SELECT name, price FROM products WHERE name LIKE '%" + q + "%'"
		log.Println("search query:", sqlText)
		rows, err := db.Query(sqlText)
		if err != nil {
			render(w, view{Q: q, Msg: "SQLエラー: " + err.Error()})
			return
		}
		defer rows.Close()
		var out [][2]string
		for rows.Next() {
			var a, b string
			rows.Scan(&a, &b)
			out = append(out, [2]string{a, b})
		}
		render(w, view{Q: q, Rows: out})
	})

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
