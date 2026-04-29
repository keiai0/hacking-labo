# crypto/pki — PKI / TLS（証明書検証とMITM・電子署名）

自前の認証局(CA)を立て、サーバ証明書を発行し、「証明書を検証するかどうか」で中間者攻撃(MITM)を防げるか／防げないかを確認する。電子署名による改ざん検知も扱う。

## 構成
| ホスト | IP | 証明書 |
|---|---|---|
| client | 10.10.11.10 | CA証明書 `ca.crt` を信頼（curl/openssl） |
| legit | 10.10.11.20 | `bank.local`：**CAが署名**した正規証明書 |
| mitm | 10.10.11.30 | `bank.local`：**自己署名**（CAの署名なし＝なりすまし） |

`pkigen` が起動時に CA・正規・なりすましの各証明書を生成し、共有ボリューム `/pki` に置く。

## 攻撃と対策の確認
```sh
make up    LAB=crypto/pki
make shell LAB=crypto/pki SVC=client
```

### 証明書検証とMITM
```sh
# 1) 正規サーバへ（CA信頼あり・検証する）→ 成功
curl --cacert /pki/ca.crt --resolve bank.local:443:10.10.11.20 https://bank.local/

# 2) なりすましサーバへ（検証する）→ 拒否（SSL verify failed, exit 60）
curl --cacert /pki/ca.crt --resolve bank.local:443:10.10.11.30 https://bank.local/

# 3) なりすましサーバへ（検証しない -k）→ MITM成立（偽サーバに接続してしまう）
curl -k --resolve bank.local:443:10.10.11.30 https://bank.local/
```
| 接続先 | 検証 | 結果 |
|---|---|---|
| 正規 | あり | 成功 |
| なりすまし | あり | **拒否**（証明書がCAで検証できない） |
| なりすまし | なし(`-k`) | **MITM成立** |

要点: TLSの安全性は「証明書の検証」で成り立つ。`-k`（検証無効）や証明書エラーの無視は、暗号化していてもMITMを許す（[[network/playground]] のTLS-MITMと対応）。

### 証明書の中身を見る
```sh
openssl x509 -in /pki/legit.crt -noout -issuer -subject   # issuer=CA, subject=bank.local
openssl x509 -in /pki/mitm.crt  -noout -issuer -subject   # issuer=自分（自己署名）
openssl s_client -connect 10.10.11.30:443 -servername bank.local </dev/null 2>/dev/null | openssl x509 -noout -issuer
```

## 電子署名（改ざん検知）
```sh
echo '振込先: 正規口座' > order.txt
openssl dgst -sha256 -sign /pki/ca.key -out order.sig order.txt          # 秘密鍵で署名
openssl x509 -in /pki/ca.crt -pubkey -noout > capub.pem
openssl dgst -sha256 -verify capub.pem -signature order.sig order.txt    # → Verified OK
echo '振込先: 攻撃者口座' > order.txt                                     # 改ざん
openssl dgst -sha256 -verify capub.pem -signature order.sig order.txt    # → Verification Failure
```
秘密鍵で署名し公開鍵で検証。内容が変わると検証に失敗する（完全性・真正性）。

## 試験での問われ方
- PKIの構成要素（CA・証明書・CRL/OCSP・信頼チェーン）
- サーバ証明書の検証（発行者・有効期限・CN/SAN・失効）、証明書検証不備
- 公開鍵暗号と電子署名、ハッシュ、完全性/真正性/否認防止
- TLSハンドシェイク、HSTS、証明書ピンニング

## 後片付け
```sh
make down LAB=crypto/pki
```
