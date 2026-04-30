# network/playground — ネットワーク演習環境（標的＋トグル）

小さな社内ネットワークを1つ用意し、偵察・盗聴・ARP・DNS・TLS-MITM などを**同じ環境**で行う。防御（FW）や監視（IDS）、キャプチャは**トグル**で後付けする。

## 構成
IPは `compose.yml` で固定している。下表の値はそこから来ている。

| ホスト | IP | 役割 |
|---|---|---|
| attacker | 10.10.0.10 | Kali（nmap, tcpdump, dsniff 等）。ここから攻撃する |
| web | 10.10.0.20 | 平文サービス群: FTP / Telnet / HTTP / HTTP Basic(/secret) / 管理画面(:8080) / SSH |
| secure | 10.10.0.30 | 暗号化サービス: SSH・SFTP / HTTPS（web との比較用） |
| victim | 10.10.0.50 | 利用者役。15秒ごとに web/secure へ**自動ログイン**し通信を流す（盗聴の的） |

- 認証情報 `alice` / `S3cretPass!` は victim がログインに使うもの（`web`/`secure`/`victim` の設定で定義）。盗聴で復元できるのがこのパスワード。
- トグル（必要時のみ起動）。各ラボ共通のMakeターゲットで、`LAB=network/playground` を付ける。

| トグル | コマンド | 内容 |
|---|---|---|
| キャプチャ | `make capture` / `make capture-stop` | ブリッジ全体を pcap 記録し `captures/` へ出力 |
| IDS | `make ids` / `make ids-stop` | Suricata で監視し `captures/suricata/eve.json` にアラート |
| ファイアウォール | `make fw-on` / `make fw-off` | web の iptables を 22,80 のみ許可にライブ切替 |

## 進め方（共通）
- `make ...` は **Mac のシェル**（このリポジトリ直下）で実行する。
- 攻撃ツール（nmap/dig/tshark 等）は **attacker コンテナの中**で実行する。入るには:
  ```sh
  make up    LAB=network/playground     # 環境を起動（最初の1回）
  make shell LAB=network/playground     # attacker コンテナに入る（以後のツールはこの中で）
  ```
- `captures/` は Mac と共有。pcap や eve.json は Mac の Wireshark/エディタでも開ける。

---

## 演習1) 偵察（recon）— 公開されているサービスを洗い出す
**目的**: 攻撃者が最初に行う偵察。開いているポート・サービス・バージョンを調べ、不要な公開がどれだけ手がかりになるかを体感する。

**手順**（attacker コンテナ内）
1. web(10.10.0.20) をポートスキャンする。nmap はポートスキャンのツール。
   ```sh
   nmap -sV -p- 10.10.0.20
   ```
   - `-p-`: 全65535ポートを対象 / `-sV`: サービスとバージョンを検出
   - → FTP(21)/SSH(22)/Telnet(23)/HTTP(80)/管理画面(8080) が見つかる。nginx のバージョンも表示され、攻撃者は既知脆弱性を調べる手がかりにする。
2. 見つかった Web を覗く。このラボでは web の nginx に**わざとディレクトリ一覧表示（`autoindex on`）を設定している**（`web/nginx.conf`。nginx の既定は `off` なので、本来は一覧が出ないのが正しい状態。ここは設定ミスを再現したもの）。そのため、**攻撃者はファイル名を知らなくても** `/backup/` の中身を一覧で見られ、置き忘れの `db_2026.sql` に辿り着ける（`autoindex` が既定どおり `off` なら、攻撃者は `ffuf`/`gobuster` 等でファイル名を総当りして探す必要がある）。
   ```sh
   curl -s http://10.10.0.20/backup/      # フォルダ内のファイル一覧が返る → db_2026.sql が見える
   curl -s http://10.10.0.20:8080/        # web が認証なしで公開している管理画面
   ```
3. **対策の確認**: ファイアウォールを有効化して再スキャンすると、許可した 22,80 以外が `filtered` になる。
   ```sh
   # Mac のシェルで:  make fw-on LAB=network/playground
   nmap -Pn -p- 10.10.0.20                # attacker 内で再スキャン → 21,23,8080 が filtered
   ```

## 演習2) 盗聴（sniffing）— 平文通信からパスワードを盗む
**目的**: victim が流す通信を傍受し、平文プロトコル(FTP/Telnet/HTTP Basic)では認証情報が丸見え、暗号化(SSH/HTTPS)では読めないことを比較する。

**前提**: 演習1-3でファイアウォールを有効化した場合は、先に戻す。有効なままだと web の FTP(21)/Telnet(23) が遮断され、平文の認証情報が流れず録れない（HTTPだけ録れて「内容が見れない」状態になる）。
```sh
make fw-off LAB=network/playground      # Mac のシェル。web の 21,23 を再び通す
```

**手順**
1. 通信を記録する（**Mac のシェル**）。ブリッジ全体を tcpdump し `captures/` に pcap を出力する。victim は15秒ごとに FTP/Telnet/HTTP/SSH/HTTPS へログインするので、数回ぶん溜める。
   ```sh
   make capture LAB=network/playground
   # 40秒ほど待つ（victim が各サービスへ数回ログインする）
   make capture-stop LAB=network/playground
   ```
   → `captures/playground_<日時>.pcap` が生成される（これが次の `$pcap`）。
2. 中身を抽出する（attacker コンテナ内。`/captures` が共有されている）。
   ```sh
   pcap=$(ls -t /captures/playground_*.pcap | head -1)   # 直近のpcapを変数に

   tshark -r "$pcap" -Y 'ftp.request.command=="PASS"' -T fields -e ftp.request.arg
   # → S3cretPass!  （FTPログインのパスワードが平文で流れている）

   tshark -r "$pcap" -Y http.authbasic -T fields -e http.authbasic
   # → alice:S3cretPass!  （HTTP Basic認証も平文）
   ```
   暗号化側(secure, ポート22/443)の通信を同様に追っても、パスワードは出てこない。
   ※ Mac の Wireshark で `captures/` の同じ pcap を開いてもよい。

## 演習3) 検知（IDS）— 防御側が平文通信を検知する
**目的**: 演習2の平文通信を、守る側が Suricata（IDS）で検知できることを確認する。

**手順**
1. IDS を起動する（**Mac のシェル**）。Suricata がブリッジを監視し、アラートを `captures/suricata/eve.json` に書く。検知ルールは `labs/network/playground/ids/custom.rules`（平文FTP/Telnet などを alert）。
   ```sh
   make ids LAB=network/playground
   ```
2. 数十秒後、アラートを集計する（**Mac のシェル**。`captures/` はMac上）。
   ```sh
   grep -o '"signature":"[^"]*"' captures/suricata/eve.json | sort | uniq -c
   # → "Cleartext FTP detected" / "Cleartext Telnet detected" が計上される
   ```

## 演習4) DNS ゾーン転送（AXFR）— 内部構成の漏洩
**目的**: 内部DNS `dns`(10.10.0.40) の設定ミス（`allow-transfer { any; }`）により、ゾーン全体（内部ホスト一覧）を誰でも取得できてしまうことを確認する。

**手順**（attacker コンテナ内）
```sh
dig +noall +answer @10.10.0.40 corp.local AXFR
```
- `@10.10.0.40`: 問い合わせ先のDNS / `corp.local`: 対象ゾーン / `AXFR`: ゾーン転送要求
- → `vpn` / `backup` / `jenkins-dev` / `admin` など、本来秘すべき内部ホストが一覧で返る

**対策**: `labs/network/playground/dns/named.conf` の `allow-transfer` をセカンダリに限定（または `{ none; }`）、再帰も内部に限定する。

## 演習5) ARPスプーフィング（中間者攻撃）
**目的**: 同一セグメントで ARP を詐称し、victim↔web の通信を自分（attacker）経由に引き込む（盗聴・改ざんの前段）。
※ 自動傍受までは行わず、手動で仕組みを体験する。

**手順**（attacker コンテナ内。`eth0` は attacker の NIC）
```sh
sysctl -w net.ipv4.ip_forward=1                     # 転送を有効化（通信を素通しさせ気づかれにくく）
arpspoof -i eth0 -t 10.10.0.50 10.10.0.20           # victim に「webは自分」と詐称（別ターミナル）
arpspoof -i eth0 -t 10.10.0.20 10.10.0.50           # web に「victimは自分」と詐称（別ターミナル）
# 経路に入った状態で、演習2の make capture / tshark で平文を取得できる
```
**対策**: スイッチの DAI（Dynamic ARP Inspection）、静的ARP、通信の暗号化。

## 演習6) TLS証明書の検証不備（HTTPSのMITM）
**目的**: 暗号化していても「証明書を検証しなければ」MITMが成立することを理解する。victim は `secure` へ `curl -k`（証明書を検証しない）でアクセスしている設定。

**手順の考え方**: 演習5で経路に入った状態で `mitmproxy` を透過モードで挟むと、victim が証明書エラーを無視するため HTTPS でも復号できてしまう。
```sh
mitmproxy --mode transparent
```
**対策**: 証明書を必ず検証する（`-k` を使わない／証明書ピンニング）、HSTS。検証していれば警告で接続が遮断される、という点が要点。より厳密な検証の実機は [[crypto/pki]] を参照。

---

## 試験での問われ方
- ポートスキャンと攻撃面の縮小、FWのデフォルトDENY
- 平文プロトコルの危険性と TLS/SSH による機密性確保
- IDS/IPS の役割（検知と防御の違い）、シグネチャ型の限界
- スイッチ環境での盗聴（ミラーポート/ARPスプーフィング）

## 後片付け
```sh
make down LAB=network/playground   # トグル類もまとめて停止・削除
```
