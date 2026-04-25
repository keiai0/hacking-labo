// 対策版: プレースホルダ（プリペアドステートメント）＋パスワードはbcryptで検証
package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var page = template.Must(template.New("p").Parse(`<!doctype html><meta charset="utf-8">
<title>Shop (fixed)</title>
<h1>社員ポータル（対策版）</h1>
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
		// 対策: プレースホルダでユーザを取得し、パスワードはハッシュで検証
		var name, hash string
		err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", u).Scan(&name, &hash)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) != nil {
			render(w, view{Msg: "ログイン失敗"})
			return
		}
		render(w, view{Msg: "ようこそ、" + name + " さん（ログイン成功）"})
	})

	http.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		// 対策: 値をプレースホルダで渡す（SQL構文として解釈されない）
		rows, err := db.Query("SELECT name, price FROM products WHERE name LIKE '%' || ? || '%'", q)
		if err != nil {
			render(w, view{Q: q, Msg: "SQLエラー"})
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
	// パスワードはbcryptハッシュで保存
	for _, u := range []struct{ name, pw string }{{"admin", "S3cr3t-admin!"}, {"alice", "alicepw"}} {
		h, _ := bcrypt.GenerateFromPassword([]byte(u.pw), bcrypt.DefaultCost)
		db.Exec(`INSERT INTO users(username,password) VALUES(?,?)`, u.name, string(h))
	}
	db.Exec(`INSERT INTO products(name,price) VALUES('Keyboard','3000'),('Mouse','1500'),('Monitor','20000')`)
}
