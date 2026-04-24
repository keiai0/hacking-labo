# network/01-recon — 偵察（ポートスキャンと情報収集）

## 目的
攻撃者が最初に行う偵察を体験する。開いているポート・稼働サービス・バージョンを洗い出し、不要な公開がどれだけ手がかりを与えるかを理解する。

## 構成
| ホスト | IP | 役割 |
|---|---|---|
| attacker | 10.10.1.10 | Kali（nmap） |
| vuln | 10.10.1.20 | 脆弱版（FTP/SSH/Telnet/HTTP/管理画面8080、ディレクトリ一覧ON、バージョン表示） |
| fixed | 10.10.1.30 | 対策版（SSH/HTTPのみ、iptablesでDROP、一覧OFF・バージョン非表示） |

## 攻撃
```sh
make up LAB=network/01-recon
make shell LAB=network/01-recon

nmap -sV -p- 10.10.1.20          # 全ポート + バージョン。21,22,23,80,8080が見える
curl -s http://10.10.1.20/backup/   # ディレクトリ一覧からdb_2026.sqlが辿れる
curl -s http://10.10.1.20:8080/     # 認証なしの管理画面
```

## 対策（fixed との差分）
- 不要サービス（FTP/Telnet/8080）を停止 → 攻撃面を減らす
- iptables で許可ポート以外を DROP（`nmap 10.10.1.30` が filtered になる）
- nginx: `autoindex off` / `server_tokens off`

```sh
nmap -Pn -p- 10.10.1.30          # 22,80 以外は filtered
```

## 試験での問われ方
- ポートスキャン手法（TCPコネクト / SYN / ステルス）、`-sV` の役割
- 攻撃面（アタックサーフェス）の縮小、不要サービスの停止
- バナー・バージョン情報の秘匿、ディレクトリリスティングの無効化
- FWのデフォルトDENY（ホワイトリスト方式）

## 後片付け
```sh
make down LAB=network/01-recon
```
