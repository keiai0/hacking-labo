# ロードマップ / 残作業

進捗と、今後の拡充候補をまとめる。`[x]`=動作確認済み / `[~]`=手順記載のみ（未自動化）/ `[ ]`=未着手。

## 構築済みの標的環境
- [x] **crypto/pki** — 自前CA、証明書検証とMITM、電子署名（改ざん検知）
- [x] **network/playground** — 偵察(nmap)、盗聴(平文↔暗号)、FWトグル、IDS(Suricata)、DNSゾーン転送
  - [~] ARPスプーフィング（ツールは用意、手順のみ）
  - [~] TLS証明書検証不備のMITM（手順のみ）
- [x] **web/webapp** — 11ページ（SQLi/XSS/OSコマンド/パストラバーサル/SSRF/リダイレクト/CORS/CSRF/セッション固定/JWT/総当たり/XXE(PHP)）＋WAFトグル
- [x] **mail** — オープンリレー、SPFによるなりすまし拒否
  - [~] DKIM署名/検証（解説のみ）、DMARC
- [x] **analysis/password-crack** — John/hashcat（raw-md5/md5crypt/zip/bcrypt）
- [x] **analysis/forensics** — pcapからC2ビーコン発見・持ち出しデータ復元
- [x] **secure-coding** — バッファオーバーフロー(認証バイパス/カナリア)、整数/符号検査
- [x] **container** — 特権 vs ハードニングの比較（観察）
  - [~] 実脱出（ディスクマウント/docker.sock）＝使い捨てVM内の手順
- [x] **scenarios/corp-network** — 外部/DMZ/内部のFW分離、到達性検証、侵入筋書き

## 残り（最後に実施予定）
- [x] **auth（OAuth/OIDC）＝軽量版**（`auth`）— redirect_uri検証不備・state欠如・PKCEなし・IDトークン検証不備の4種を自作IdP+RPで確認済み。
  - 実機Keycloakで特定設定を触りたくなったら別途追加。

当初計画の全分野を構築完了。以降は下記の拡充候補のみ。

## 拡充候補（任意・優先度順の目安）

### 高：試験頻出・手を動かす価値が高い
- [x] **PKI/TLS 専用ラボ**（`crypto/pki`）— 自前CA、サーバ証明書発行、検証あり/なしでMITMの成否、電子署名の改ざん検知。
- [ ] **検知基盤（軽量SIEM）トグル** — rsyslog→Loki→Grafana でログ集約・相関。各環境に後付けするトグル化。
- [ ] **webapp 追加ページ** — ファイルアップロード検証不備、IDOR、SSTI、LFI/RFI(PHP)、セキュリティヘッダ/CSP、WAFバイパス演習。

### 中：分野の穴埋め
- [ ] **secure-coding 追加** — 競合状態(TOCTOU)、フォーマット文字列、use-after-free。
- [ ] **クラウド設定ミス（container内）** — MinIOの公開バケット匿名アクセス、メタデータSSRF（web appと連携）。
- [ ] **analysis 追加** — メモリフォレンジック(Volatility 3)、ディスク(The Sleuth Kit)。学習用サンプルイメージが必要。
- [ ] **ログ解析ラボ** — Web/認証ログから侵入経路をたどる（forensicsの延長）。

### 低：発展・後回し
- [ ] **playground 追加** — VPN(WireGuard)、nftablesでの本格FW設計、DNSキャッシュポイズニングの仕組み。
- [ ] **scenario 作り込み** — DMZにメール、内部に業務PC、インシデント対応ワークシート。
- [ ] **Windows / Active Directory** — Kerberos, Pass-the-Hash, Kerberoasting（保留。実施時はUTM等のVM）。

## 横断的な整理
- [ ] **docs/** の充実 — 用語集、関連法規、午後問題の振り返りテンプレート。
- [ ] イメージのバージョン固定の点検（`nginx:1.27` 等のタグ）。
- [ ] 各ラボREADMEの「試験での問われ方」を過去問と突き合わせて加筆。

## メモ
- 自動化しない方針のもの: ARPスプーフィング等の実傍受、コンテナ実脱出（安全上、手順記載＋各自実施/VM内）。
- Makefileは `labs/` 配下向け。`scenarios/` は `docker compose -p hl-scenario` を直接使う。
