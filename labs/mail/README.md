# mail — メールのなりすまし対策（SPF/DKIM/DMARC）とオープンリレー

送信ドメイン認証がないと、第三者が正規ドメインを騙ってメールを送れる。また中継制限がないと踏み台（オープンリレー）にされる。脆弱版/対策版のメールサーバを並べて比較する。

## 構成
| ホスト | IP | 役割 |
|---|---|---|
| attacker | 10.10.6.10 | swaks で送信（partner.example の SPF に**載っていない**IP） |
| dns | 10.10.6.40 | partner.example の SPF/DMARC レコードを提供 |
| mailv | 10.10.6.20 | 脆弱版: オープンリレー + SPF/DMARC検証なし |
| mailf | 10.10.6.30 | 対策版: 中継制限 + SPF検証（policyd-spf） |

`partner.example` のDNS: `v=spf1 ip4:10.10.6.99 -all`（正規送信元は .99 のみ）/ `_dmarc` は `p=reject`。

## 攻撃と対策の確認
```sh
make up LAB=mail
make shell LAB=mail
```

### 1) オープンリレー（踏み台化）
外部ドメイン宛の中継を依頼する。
```sh
swaks --server 10.10.6.20 --from x@attacker.test --to victim@gmail.com    # mailv → 250 Ok（中継してしまう）
swaks --server 10.10.6.30 --from x@attacker.test --to victim@gmail.com    # mailf → 454 Relay access denied
```

### 2) 送信ドメインなりすまし（SPF）
SPFに載らない 10.10.6.10 から `From: ceo@partner.example` を送る。
```sh
swaks --server 10.10.6.20 --from ceo@partner.example --to alice@corp.local   # mailv → 250 Ok（受理）
swaks --server 10.10.6.30 --from ceo@partner.example --to alice@corp.local   # mailf → 550 SPF fail - not authorized
```
mailf は送信元IPを partner.example のSPFと照合し、不一致を拒否する。DMARC（`p=reject`）はSPF/DKIMの失敗を「拒否」に結びつけるポリシーで、本ラボではSPF失敗がそのまま拒否に対応する。

## 補足（DKIM）
DKIMは本文・ヘッダへの電子署名で、改ざん検知と送信元認証を行う。鍵配布（DNSのTXT）と OpenDKIM の署名・検証が必要で本ラボでは省略。試験では「SPF=IP認証、DKIM=署名、DMARC=両者の整合と処理方針」の役割分担を押さえる。

## 試験での問われ方
- SPF / DKIM / DMARC それぞれの仕組みと防げる脅威の違い
- オープンリレーの危険性（スパム踏み台）、中継制限（`reject_unauth_destination`）
- エンベロープFrom と ヘッダFrom の違い、DMARCアライメント
- なりすましメール・BEC（ビジネスメール詐欺）対策

## 後片付け
```sh
make down LAB=mail
```
