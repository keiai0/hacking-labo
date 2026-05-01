# analysis/forensics — ネットワークフォレンジック（C2通信の発見）

感染ホストが外部のC2サーバへ定期的に通信（ビーコン）し、情報を持ち出す様子をキャプチャし、pcapから不審通信を特定・解析する。

## 構成
| ホスト | IP | 役割 |
|---|---|---|
| analysis | 10.10.7.10 | 解析役（tshark 等）。Macの Wireshark でも可 |
| infected | 10.10.7.20 | 感染ホスト。5秒ごとにC2へビーコン |
| c2 | 10.10.7.66 | 攻撃者のC2サーバ（`/api/beacon` に200を返す） |
| web | 10.10.7.30 | 正規サイト（ノイズ用） |

## 手順
```sh
make up      LAB=analysis/forensics
make capture LAB=analysis/forensics   # 記録開始
# 40〜60秒ほど待つ（ビーコンが数回発生）
make capture-stop LAB=analysis/forensics
make shell   LAB=analysis/forensics
```
pcap は `captures/forensics/*.pcap`（Mac の Wireshark でも開ける）。

## 解析のポイント
**目的**: 解析役（防御側）が、正規通信に紛れた感染ホスト→C2のビーコンをpcapから見分け、何が持ち出されたかを復元する。手がかりを知らない状態から「宛先の偏り → 一定周期 → 不審なUA/URI → 中身の復元」と辿るのがこの演習。

**このラボが流しているもの**（`infected/beacon.sh`）: 感染ホスト(10.10.7.20)が5秒ごとに C2(10.10.7.66) の `/api/beacon` へPOSTし、持ち出しデータ `user=alice;card=4111111111111111;host=INFECTED01` を**base64化してフォーム項目 `d=` に詰めて**送る。User-Agentは `beacon/1.0`。合間に正規サイト(10.10.7.30)へのアクセスをノイズとして混ぜている。下の手順の `ip.dst==10.10.7.66` や `d=` はこの通信に由来する。
```sh
f=$(ls -t /captures/forensics/*.pcap | head -1)

# 1) 宛先別のパケット数 → 特定ホストへの反復通信が目立つ
tshark -r "$f" -q -z ip_dsts,tree

# 2) C2宛HTTPの時刻 → ほぼ一定間隔（5秒周期）= ビーコンの特徴
tshark -r "$f" -Y 'http.request && ip.dst==10.10.7.66' \
  -T fields -e frame.time_relative -e http.request.method -e http.request.uri

# 3) 不審な User-Agent / URI（/api/beacon, beacon/1.0）

# 4) 持ち出しデータの復元: C2宛POSTのフォーム値を取り出し、d= を外してbase64デコード
tshark -r "$f" -Y 'http.request && ip.dst==10.10.7.66' \
  -T fields -e urlencoded-form.value | head -1 | sed 's/^d=//' | base64 -d
# → user=alice;card=4111111111111111;host=INFECTED01（持ち出された中身が平文で復元できる）
```

ビーコンの典型的特徴: **一定間隔・同一宛先・小さく規則的なサイズ・不審なUA/URI**。正規通信（web宛）と混在する中から見分ける。

## 発展（任意）
- **メモリフォレンジック（Volatility 3）**／**ディスク（The Sleuth Kit）** は解析対象のイメージが必要。CTF等の学習用サンプルを `captures/` に置いて解析する（数GB、ここでは同梱しない）。
- ログ解析: Webサーバ/認証ログから侵入経路をたどる演習は scenario 環境と連携予定。

## 試験での問われ方
- インシデント対応のフロー（検知・封じ込め・根絶・復旧）
- C2通信・ビーコンの特徴、ふるまい検知、IoC（Indicator of Compromise）
- ログ・pcapからの侵害範囲特定、タイムライン作成
- 情報持ち出し（exfiltration）の手口と検知

## 後片付け
```sh
make down LAB=analysis/forensics
```
