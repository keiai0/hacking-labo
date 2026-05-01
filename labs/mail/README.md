# mail — メールのなりすまし対策（SPF/DKIM/DMARC）とオープンリレー

送信ドメイン認証がないと、攻撃者が正規ドメインを騙ってメールを送れる。また中継制限がないと、メールサーバが第三者のスパム中継（オープンリレー＝踏み台）に悪用される。このラボは**脆弱版と対策版のメールサーバを並べて**、同じ攻撃が片方では通り、片方では拒否されることを比較する。

## 構成
IP・ホスト名・DNSレコードはすべて `compose.yml` と `dns/db.partner.example` で定義している。下表・下記の値はそこから来ている。

| ホスト | IP | 役割 |
|---|---|---|
| attacker | 10.10.6.10 | `swaks` でメールを送る攻撃者。**partner.example のSPFに載っていないIP**（後述） |
| dns | 10.10.6.40 | `partner.example` の SPF / DMARC / MX レコードを提供するDNS |
| mailv | 10.10.6.20 | **脆弱版**メールサーバ: オープンリレー + SPF/DMARC検証なし |
| mailf | 10.10.6.30 | **対策版**メールサーバ: 中継制限 + SPF検証（policyd-spf） |

**このラボでわざと作り込んだ設定（`dns/db.partner.example`）**:
- `mail.partner.example` の正規メールサーバは **10.10.6.99** だけ。
- SPFレコード `v=spf1 ip4:10.10.6.99 -all` … 「partner.example を名乗ってよい送信元は 10.10.6.99 のみ、それ以外は拒否(`-all`=hardfail)」という宣言。
- → 攻撃者の **10.10.6.10 はこの .99 に該当しない**ので、SPFを検証するサーバ(mailf)は「正規送信元ではない」と判定できる。

**脆弱版(mailv)の中身**（`mailv/entrypoint.sh`）: Postfixの `smtpd_recipient_restrictions=permit` にしてある。これは「誰から誰宛でも中継を許可」という意味で、**わざとオープンリレー＆SPF無検証**にした状態。対策版(mailf)は `permit_mynetworks, reject_unauth_destination, check_policy_service ...policyd-spf`（自ネットワーク以外の外部宛中継を拒否 + SPF検証）。

## 進め方
```sh
make up LAB=mail
make shell LAB=mail        # 攻撃者(attacker)コンテナに入る。以降の swaks はこの中で実行
```
`swaks` は SMTP の送信を手で試すための定番ツール。`--server` で接続先メールサーバ、`--from` でエンベロープFrom（送信者）、`--to` で宛先を指定する。**同じ送信を mailv(脆弱) と mailf(対策) の両方に投げて、応答コードの違いを見る**のがこのラボの基本。

---

## 演習1) オープンリレー（踏み台化）
**目的**: 中継制限のないメールサーバは、無関係な第三者（攻撃者）に「外部ドメイン宛のメールを中継して」と頼まれても引き受けてしまい、スパムの踏み台になることを確認する。

**手順**（attacker コンテナ内）。attacker 自身とは無関係な `gmail.com` 宛の中継を、mailv と mailf に依頼する。
```sh
swaks --server 10.10.6.20 --from x@attacker.test --to victim@gmail.com   # → mailv
swaks --server 10.10.6.30 --from x@attacker.test --to victim@gmail.com   # → mailf
```
- **mailv(脆弱)** は `250 Ok` を返す。自分の利用者でも自分宛でもない外部→外部のメールを**中継してしまう**（＝オープンリレー。攻撃者はこのサーバを使って大量スパムを送れる）。
- **mailf(対策)** は `454 Relay access denied` で拒否する。`reject_unauth_destination` により、自分が受け持つドメイン以外への中継を断っているため。

## 演習2) 送信ドメインなりすまし（SPF）
**目的**: SPFを検証しないサーバは、攻撃者が `From: ceo@partner.example` と名乗ったメールを本物として受理してしまう（なりすまし／BECの入口）。SPF検証があれば送信元IPの不一致で弾けることを確認する。

**手順**（attacker コンテナ内）。攻撃者(10.10.6.10)が、**partner.example のSPFに載っていないIPから** `ceo@partner.example` を騙って社内の `alice` 宛に送る。
```sh
swaks --server 10.10.6.20 --from ceo@partner.example --to alice@corp.local   # → mailv
swaks --server 10.10.6.30 --from ceo@partner.example --to alice@corp.local   # → mailf
```
- **mailv(脆弱)** は `250 Ok`。送信元IPとSPFレコードを照合しないので、**攻撃者が誰を騙っても受理**してしまう。
- **mailf(対策)** は `550 SPF fail - not authorized`。mailf は送信元IP(10.10.6.10)を partner.example のSPF(`ip4:10.10.6.99`)と照合し、不一致なので拒否する。

DMARC（`_dmarc` の `p=reject`）は、SPF/DKIMの検証失敗を「どう扱うか（監視/隔離/拒否）」を送信ドメイン側が宣言するポリシー。このラボでは `p=reject` なので、SPF失敗がそのまま「拒否」に対応する。

## 補足（DKIM）
DKIMは本文・ヘッダへの電子署名で、改ざん検知と送信元認証を行う（SPFはIPの照合、DKIMは署名の照合）。鍵配布（DNSのTXTに公開鍵）と OpenDKIM による署名・検証の構築が必要で、このラボでは省略している。試験では「**SPF=送信元IPの認証 / DKIM=署名による認証と改ざん検知 / DMARC=SPF・DKIMの整合(アライメント)確認と処理方針**」の役割分担を押さえる。

## 試験での問われ方
- SPF / DKIM / DMARC それぞれの仕組みと防げる脅威の違い
- オープンリレーの危険性（スパム踏み台）、中継制限（`reject_unauth_destination`）
- エンベロープFrom と ヘッダFrom の違い、DMARCアライメント
- なりすましメール・BEC（ビジネスメール詐欺）対策

## 後片付け
```sh
make down LAB=mail
```
