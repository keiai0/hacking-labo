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
| `/cmdi` | OSコマンドインジェクション | `sh -c` に入力連結 | シェルを介さず引数渡し |
| `/path` | パストラバーサル | ベース+入力を連結 | `filepath.Base`＋配下チェック |
| `/ssrf` | SSRF | 任意URLをサーバ取得 | スキーム/ホストをallowlist |
| `/redirect` | オープンリダイレクト | 入力URLへそのまま遷移 | 相対パスのみ許可 |
| `/cors` | CORS設定ミス | Originを無検証で反射 | allowlist一致のみ許可 |
| `/csrf` | CSRF | トークン検証なし | トークン必須（不一致は403） |
| `/session` | セッション固定 | ログイン時にID再生成せず | ログイン時に新ID発行 |
| `/jwt` | JWT検証不備 | alg=none を受理 | 署名必須・alg=none拒否 |
| `/brute` | ログイン総当たり | 試行無制限 | 5回失敗でロック(429) |
| `/xxe` | XXE（外部実体） | libxmlで実体展開（PHP） | 展開しない |

`/xxe` のみ PHP サイドカー（`phpapp`）で動作し、Goアプリが `/xxe` をプロキシする（Go標準XMLは外部実体を扱わないため）。

## 動作確認済みの代表入力
| ページ | 入力 | 脆弱版 | 対策版 |
|---|---|---|---|
| /sqli | `admin' -- ` | ログイン成功 | 失敗 |
| /cmdi | `x; id` | idが実行 | 文字列のまま |
| /path | `../../../../etc/passwd` | 読める | エラー |
| /ssrf | `http://app:8080/sqli` | 内部へ到達 | 拒否 |
| /redirect | `?to=https://example.com/` | 外部へ遷移 | 400 |
| /cors | `Origin: https://evil.example` | 反射 | 反射しない |
| /csrf | トークン無しPOST | 変更成功 | 403 |
| /session | 事前SIDでログイン | SIDが有効化 | 新SID発行・無効 |
| /jwt | alg=none & role=admin | 管理者許可 | 拒否 |
| /brute | 連続失敗 | 無制限(200) | 5回でロック(429) |
| /xxe | 外部実体で/etc/passwd | 展開される | 展開されない |

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
