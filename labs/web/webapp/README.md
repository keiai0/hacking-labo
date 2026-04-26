# web/webapp — 統合脆弱Webアプリ

Web系の脆弱性を1つのGoアプリにまとめた演習環境。**1ページ＝1脆弱性**で、各ページ `?mode=vuln|fixed` で脆弱版/対策版を切り替える。WAFはトグルで前段に追加する。

## 構成
| アクセス先 | 経路 | WAF |
|---|---|---|
| http://127.0.0.1:8080 | ブラウザ → アプリ直結 | なし |
| http://127.0.0.1:8090 | ブラウザ → ModSecurity(CRS) → アプリ | あり（`make waf-on` 時） |

ブラウザ（Firefox）＋ Burp Suite は Mac から操作。攻撃用コンテナは `make shell`。

## 実装済みページ
| パス | 脆弱性 | 脆弱版の実装 | 対策版の実装 |
|---|---|---|---|
| `/sqli` | SQLインジェクション | 文字列連結でSQL生成 | プレースホルダ |
| `/xss` | 反射型XSS | 入力を生でHTML出力 | HTMLエスケープ |

（CSRF, SSRF, OSコマンド, パストラバーサル, XXE, リダイレクト, CORS, セッション, JWT … を同じ要領で追加していく）

## 使い方
```sh
make up LAB=web/webapp
```
ブラウザで http://127.0.0.1:8080 → メニューからページ選択。各ページの「対策版に切替」で挙動を比較。

代表的な攻撃入力:
- SQLi ログイン: ユーザ名 `admin' -- `
- SQLi 検索（UNION）: `x' UNION SELECT username,password FROM users -- `
- XSS: 名前欄に `<script>alert(1)</script>`

sqlmap（攻撃用コンテナ）:
```sh
make shell LAB=web/webapp
sqlmap -u 'http://app:8080/sqli?q=1' --batch --dump -T users
```

## WAFトグル
```sh
make waf-on  LAB=web/webapp      # 前段にModSecurity+OWASP CRSを起動
# → :8080 の攻撃は成功、:8090 経由は 403 でブロックされる
make waf-off LAB=web/webapp
```
WAFをすり抜けるエンコーディングも試し、「WAFは多層防御の補助であって根本対策ではない」ことを確認する。

## 試験での問われ方
- 各脆弱性の原理と根本対策（入力検証・出力エスケープ・プレースホルダ・トークン等）
- WAFの役割と限界、多層防御
- セッション管理、認証・認可（ページ追加に応じて）

## 後片付け
```sh
make down LAB=web/webapp   # WAF等のトグルもまとめて停止
```
