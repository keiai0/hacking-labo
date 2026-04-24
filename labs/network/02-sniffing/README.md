# network/02-sniffing — 盗聴（平文プロトコルの危険性）

## 目的
平文プロトコル（FTP/Telnet/HTTP Basic認証）の通信を盗聴し、認証情報や通信内容がそのまま読めることを確認する。暗号化（SSH/SFTP・HTTPS）との差を、同じキャプチャの中で比べる。

## 構成
| ホスト | IP | 役割 |
|---|---|---|
| attacker | 10.10.2.10 | Kali（tcpdump/tshark） |
| vuln | 10.10.2.20 | 平文（FTP / Telnet / HTTP Basic） |
| fixed | 10.10.2.30 | 暗号化（SFTP / HTTPS） |
| victim | 10.10.2.40 | 15秒ごとに両方へログインする利用者役 |

## 攻撃（キャプチャ → Macで解析）
```sh
make up LAB=network/02-sniffing
make capture LAB=network/02-sniffing   # ラボ全体の通信を captures/ に記録
# 30秒ほど待つ
make capture-stop LAB=network/02-sniffing
```
`captures/network-02-sniffing_*.pcap` を Mac の Wireshark で開く。

- 平文の抽出例（攻撃用コンテナ内 `make shell` でも可）:
  ```sh
  tshark -r <pcap> -Y 'ftp.request.command=="PASS"' -T fields -e ftp.request.arg
  # → S3cretPass! が平文で見える。Telnet/HTTP Basic も同様
  ```
- 暗号化側（22/443）はペイロードを追っても `S3cretPass!` は出てこない。

## 対策
- FTP→SFTP、Telnet→SSH、HTTP→HTTPS。認証情報・通信内容をTLS/SSHで暗号化する。
- Wiresharkでは「暗号化されていると中身が読めない」ことを目で確認するのが要点。

## 試験での問われ方
- 平文プロトコルのリスク、スイッチングHub環境でのミラーポート/ARPスプーフィングによる盗聴
- TLS/SSHによる機密性の確保、HSTS、Basic認証をHTTPS前提にする理由
- 可用性・完全性・機密性のうち「機密性」が損なわれる例

## 後片付け
```sh
make down LAB=network/02-sniffing
```
