# network/playground — ネットワーク演習環境（標的＋トグル）

小さな社内ネットワークを1つ用意し、偵察・盗聴・ARP・DNS・TLS-MITM などを**同じ環境**で行う。防御（FW）や監視（IDS）、キャプチャは**トグル**で後付けする。

## 構成
| ホスト | IP | 役割 |
|---|---|---|
| attacker | 10.10.0.10 | Kali（nmap, tcpdump, dsniff 等） |
| web | 10.10.0.20 | 平文サービス群: FTP / Telnet / HTTP / HTTP Basic(/secret) / 管理画面(:8080) / SSH |
| secure | 10.10.0.30 | 暗号化サービス: SSH・SFTP / HTTPS（web との比較用） |
| victim | 10.10.0.50 | 利用者役。15秒ごとに web/secure へログインし通信を流す |

トグル（必要時のみ起動）:
| トグル | コマンド | 内容 |
|---|---|---|
| キャプチャ | `make capture` / `capture-stop` | ブリッジ全体を pcap 記録（→ Mac の Wireshark） |
| IDS | `make ids` / `ids-stop` | Suricata で監視し captures/suricata/eve.json にアラート |
| ファイアウォール | `make fw-on` / `fw-off` | web の iptables を 22,80 のみ許可にライブ切替 |

いずれも `LAB=network/playground` を付ける。

## 演習

### 1) 偵察（recon）
```sh
make up LAB=network/playground
make shell LAB=network/playground
nmap -sV -p- 10.10.0.20        # FTP/Telnet/HTTP/8080 が見える
curl -s http://10.10.0.20/backup/     # ディレクトリ一覧 → db_2026.sql
curl -s http://10.10.0.20:8080/       # 認証なし管理画面
```
対策確認: `make fw-on` 後に再スキャン → 22,80 以外が filtered に。

### 2) 盗聴（sniffing）
```sh
make capture LAB=network/playground   # 30秒ほど待つ
make capture-stop LAB=network/playground
```
pcap を Wireshark で開く／tshark で抽出:
```sh
tshark -r <pcap> -Y 'ftp.request.command=="PASS"' -T fields -e ftp.request.arg   # S3cretPass!
tshark -r <pcap> -Y http.authbasic -T fields -e http.authbasic                   # alice:S3cretPass!
```
→ 平文(web)は筒抜け。暗号化(secure, 22/443)は中身が読めない。

### 3) 検知（detection / IDS）
```sh
make ids LAB=network/playground       # Suricata起動（別ターミナルで放置）
# victim の平文通信にアラートが上がる
grep -o '"signature":"[^"]*"' captures/suricata/eve.json | sort | uniq -c
```
→ Cleartext FTP / Telnet を検知。ルールは `ids/custom.rules`。

### 4) ARP / DNS / TLS-MITM（追加予定）
同じ環境に手順を追加していく。ARPスプーフィングは attacker に dsniff/arpspoof を用意済み（手順はこの環境で実施）。

## 試験での問われ方
- ポートスキャンと攻撃面の縮小、FWのデフォルトDENY
- 平文プロトコルの危険性と TLS/SSH による機密性確保
- IDS/IPS の役割（検知と防御の違い）、シグネチャ型の限界
- スイッチ環境での盗聴（ミラーポート/ARPスプーフィング）

## 後片付け
```sh
make down LAB=network/playground   # トグル類もまとめて停止・削除
```
