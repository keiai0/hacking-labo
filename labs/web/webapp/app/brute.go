package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	bruteMu   sync.Mutex
	failCount = map[string]int{} // username -> 連続失敗回数
)

const lockThreshold = 5

// ログイン総当たり。失敗回数の制限（ロック）があるかどうか。
func bruteHandler(w http.ResponseWriter, r *http.Request) {
	m := mode(r)
	var b strings.Builder

	if r.Method == http.MethodPost {
		r.ParseForm()
		u, p := r.FormValue("u"), r.FormValue("p")

		bruteMu.Lock()
		locked := m == "fixed" && failCount[u] >= lockThreshold
		bruteMu.Unlock()

		if locked {
			// 対策: しきい値超過でロック
			http.Error(w, "アカウントロック中（試行回数超過）", http.StatusTooManyRequests)
			return
		}
		if authUser(u, p) {
			bruteMu.Lock()
			failCount[u] = 0
			bruteMu.Unlock()
			b.WriteString(`<p><b>ログイン成功</b></p>`)
		} else {
			bruteMu.Lock()
			failCount[u]++
			n := failCount[u]
			bruteMu.Unlock()
			fmt.Fprintf(&b, `<p>ログイン失敗（連続%d回）</p>`, n)
		}
	}

	fmt.Fprintf(&b, `<h2>ログイン</h2><form method=post action="/brute?mode=%s">
ユーザ名 <input name=u value="alice"> パスワード <input name=p> <button>login</button></form>`, m)
	b.WriteString(`<hr><p style="color:#666">攻撃: <code>hydra</code> 等で総当たり。
脆弱版は何回でも試せる／対策版は5回失敗でロック(429)。復帰は down/up。</p>`)
	layout(w, "/brute", "ログイン総当たり", m, b.String())
}
