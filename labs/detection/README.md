# detection — ログ集約と検知（軽量SIEM: Loki + Grafana）

認証ログを一元集約し、LogQLで不審なふるまい（総当たりの兆候）を検知・可視化する。Elasticsearch系より軽量な Loki + Grafana + Promtail 構成（16GBでも可）。

## 構成
| コンテナ | 役割 | 公開 |
|---|---|---|
| loki | ログ蓄積・LogQLクエリ | 127.0.0.1:3100 (API) |
| promtail | `/logs/*.log` を収集して Loki へ転送 | — |
| grafana | 可視化（Lokiデータソースを provisioning 済み） | 127.0.0.1:3000 |
| loggen | 認証ログを生成（特定IPからの連続失敗＝総当たりを含む） | — |

## 使い方
```sh
make up LAB=detection
# ブラウザ: http://127.0.0.1:3000 （匿名Admin）→ Explore → Lokiを選択
```

## 検知（LogQL）
```sh
# 1) 失敗ログインの抽出
curl -s --get http://127.0.0.1:3100/loki/api/v1/query_range \
  --data-urlencode 'query={job="auth"} |= "Failed password"' \
  --data-urlencode "start=$(($(date +%s)-300))000000000" \
  --data-urlencode "end=$(date +%s)000000000" | head

# 2) 直近1分の失敗数（しきい値超過＝総当たり。アラートの基礎）
curl -s --get http://127.0.0.1:3100/loki/api/v1/query \
  --data-urlencode 'query=sum(count_over_time({job="auth"} |= "Failed password" [1m]))'
```
Grafana の Explore で同じクエリをグラフ表示できる。攻撃者IP `203.0.113.9` からの `Failed password` が多数（1分あたり数十件）、正規は `Accepted password` が少数、という差で総当たりを判別する。

### Grafanaでの探索例（LogQL）
- `{job="auth"}` … 全ログ
- `{job="auth"} |= "Failed password"` … 失敗のみ
- `sum by (job) (count_over_time({job="auth"} |= "Failed password" [1m]))` … 失敗数の推移（しきい値でアラート設定可）

## 他ラボとの連携
Promtail の収集対象（`promtail-config.yaml` の `__path__`）を各ラボの `captures/` やサーバログに向ければ、playground/webapp/mail/scenario のログも集約できる（＝「検知」トグルのログ版）。

## 試験での問われ方
- ログの一元管理・保全、相関分析、SIEMの役割
- 検知（しきい値/シグネチャ/ふるまい）とアラート、相関ルール
- 総当たり・不正アクセスの痕跡、ログの改ざん対策・タイムスタンプ
- 監視運用（SOC）、インシデントの早期検知

## 後片付け
```sh
make down LAB=detection
```
