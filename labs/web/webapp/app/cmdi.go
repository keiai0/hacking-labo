package main

import (
	"fmt"
	"html"
	"net/http"
	"os/exec"
	"strings"
)

// OSコマンドインジェクション。シェル経由で入力を渡すか、渡さないか。
func cmdiHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	q := r.URL.Query().Get("q")
	var b strings.Builder
	fmt.Fprintf(&b, `<h2>文字列エコー</h2><form><input type=hidden name=mode value=%s>
<input name=q value="%s"><button>実行</button></form>`, m, html.EscapeString(q))

	if q != "" {
		var out []byte
		if m == "vuln" {
			// 脆弱: シェルに入力を連結して渡す → ; や | で任意コマンド実行
			out, _ = exec.Command("sh", "-c", "echo "+q).CombinedOutput()
		} else {
			// 対策: シェルを介さず、引数として渡す（メタ文字が解釈されない）
			out, _ = exec.Command("echo", q).CombinedOutput()
		}
		fmt.Fprintf(&b, `<pre>%s</pre>`, html.EscapeString(string(out)))
	}
	b.WriteString(`<hr><p style="color:#666">試:　<code>x; id</code>（脆弱版ではidが実行される）</p>`)
	layout(w, "/cmdi", "OSコマンドインジェクション", m, b.String())
}
