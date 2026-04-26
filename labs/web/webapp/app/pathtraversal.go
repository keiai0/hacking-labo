package main

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const docRoot = "/srv/docs"

// パストラバーサル。ベースと入力を素朴に連結すると ../ で外に出られる。
func pathHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	name := r.URL.Query().Get("file")
	var b strings.Builder
	fmt.Fprintf(&b, `<h2>文書ビューア</h2><form><input type=hidden name=mode value=%s>
ファイル名 <input name=file value="%s" size=40><button>開く</button></form>`, m, html.EscapeString(name))

	if name != "" {
		var data []byte
		var err error
		if m == "vuln" {
			// 脆弱: ベース + 入力 を連結（../ を防げない）
			data, err = os.ReadFile(docRoot + "/" + name)
		} else {
			// 対策: パス名のみ採用し、docRoot配下に収まるか検証
			clean := filepath.Join(docRoot, filepath.Base(name))
			if !strings.HasPrefix(clean, docRoot+"/") {
				err = fmt.Errorf("不正なパス")
			} else {
				data, err = os.ReadFile(clean)
			}
		}
		if err != nil {
			fmt.Fprintf(&b, `<p>エラー: %s</p>`, html.EscapeString(err.Error()))
		} else {
			fmt.Fprintf(&b, `<pre>%s</pre>`, html.EscapeString(string(data)))
		}
	}
	b.WriteString(`<hr><p style="color:#666">試:　<code>../../../../etc/passwd</code>（脆弱版では読めてしまう）</p>`)
	layout(w, "/path", "パストラバーサル", m, b.String())
}
