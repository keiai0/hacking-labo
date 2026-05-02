# crypto/pki — PKI / TLS（証明書検証とMITM・電子署名）

自前の認証局(CA)を立て、サーバ証明書を発行し、「クライアントが証明書を検証するかどうか」で中間者攻撃(MITM)を防げるか／防げないかを確認する。電子署名による改ざん検知も扱う。

## 構成
証明書は `pkigen` コンテナが起動時に生成し、共有ボリューム `/pki` に置く（各コンテナから `/pki/ca.crt` などで参照できる）。

| ホスト | IP | 提示する証明書 |
|---|---|---|
| client | 10.10.11.10 | 攻撃対象ではなく利用者役。CA証明書 `ca.crt` を信頼し、curl/openssl で接続する |
| legit | 10.10.11.20 | `bank.local` の**正規**サーバ。証明書は **CAが署名**したもの |
| mitm | 10.10.11.30 | `bank.local` を騙る**なりすまし**サーバ。証明書は**自己署名**（CAの署名なし） |

**このラボで作り込んだ前提**（`pkigen/gen.sh`）:
- CA（`CN=HackLabo Root CA`）を自前で作り、client はその `ca.crt` を信頼している。
- legit の証明書と mitm の証明書は、**どちらも `CN=bank.local`（名乗りは同じ）**。違いは発行者(issuer)だけ：
  - legit … issuer = HackLabo Root CA（client が信頼するCA）→ 検証が**通る**
  - mitm … issuer = 自分自身（自己署名）→ client の信頼するCAまで辿れず検証が**失敗する**
- つまり「正しいサーバか」はドメイン名の一致だけでは決まらず、**信頼するCAが署名しているか**で決まる、というのがこのラボの核心。

## 進め方
```sh
make up    LAB=crypto/pki
make shell LAB=crypto/pki SVC=client     # 利用者(client)コンテナに入る
```

---

## 演習1) 証明書検証とMITM
**目的**: DNSを乗っ取られ `bank.local` が攻撃者(mitm)に向いても、クライアントが証明書を正しく検証していれば「このサーバは bank.local を名乗っているがCAの署名がない」と見抜いて接続を拒否できること、逆に検証を無効化(`-k`)すると暗号化していてもMITMが成立してしまうことを確認する。

**手順**（client コンテナ内）。`--resolve bank.local:443:<IP>` は「`bank.local` という名前のまま、接続先IPだけを指定したサーバに向ける」オプション。**DNSが書き換えられて bank.local が別IPを指している状況を手元で再現する**ために使う（名前は bank.local のままなので、証明書のホスト名検証はちゃんと働く）。
```sh
# 1) 正規サーバ(10.10.11.20)へ。CAを信頼して検証する → 成功
curl --cacert /pki/ca.crt --resolve bank.local:443:10.10.11.20 https://bank.local/

# 2) なりすましサーバ(10.10.11.30)へ。検証する → 拒否(SSL certificate verify failed, exit 60)
curl --cacert /pki/ca.crt --resolve bank.local:443:10.10.11.30 https://bank.local/

# 3) なりすましサーバへ。検証しない(-k) → MITM成立（client は偽サーバに気づかず接続）
curl -k --resolve bank.local:443:10.10.11.30 https://bank.local/
```
| 接続先 | クライアントの検証 | 結果 | なぜ |
|---|---|---|---|
| legit(正規) | あり(`--cacert`) | 成功 | 証明書がCAで検証でき、ホスト名も一致 |
| mitm(なりすまし) | あり(`--cacert`) | **拒否** | mitmの証明書は自己署名で、信頼するCAまで辿れない |
| mitm(なりすまし) | なし(`-k`) | **MITM成立** | 検証をやめたので client は偽証明書をそのまま受け入れる |

**要点**: TLSの安全性は暗号化そのものではなく「**証明書の検証**」で成り立つ。`-k`（検証無効）や「証明書エラーを無視」は、暗号化していてもMITMを許す。[[network/playground]] 演習6（victim が `curl -k` でアクセスしていた）のTLS-MITMと同じ話。

**証明書の中身を見て、違いを確認する**:
```sh
openssl x509 -in /pki/legit.crt -noout -issuer -subject   # issuer=HackLabo Root CA, subject=bank.local
openssl x509 -in /pki/mitm.crt  -noout -issuer -subject   # issuer=bank.local（＝自分。自己署名）, subject=bank.local
# 実際にmitmサーバが提示してくる証明書の発行者を確認
openssl s_client -connect 10.10.11.30:443 -servername bank.local </dev/null 2>/dev/null | openssl x509 -noout -issuer
```
subject（名乗り）はどちらも bank.local で同じだが、issuer（誰が署名したか）が違うことが、拒否できた理由。

## 演習2) 電子署名（改ざん検知）
**目的**: 秘密鍵で署名したデータは、公開鍵で「本人が署名したか」「改ざんされていないか」を検証できること、1文字でも内容が変わると検証に失敗することを確認する（完全性・真正性）。ここでは題材としてCAの鍵ペアを流用する。

**手順**（client コンテナ内）
```sh
echo '振込先: 正規口座' > order.txt
openssl dgst -sha256 -sign /pki/ca.key -out order.sig order.txt          # 秘密鍵(ca.key)で署名を作る
openssl x509 -in /pki/ca.crt -pubkey -noout > capub.pem                  # 証明書から公開鍵を取り出す
openssl dgst -sha256 -verify capub.pem -signature order.sig order.txt    # → Verified OK

echo '振込先: 攻撃者口座' > order.txt                                     # 攻撃者が中身を書き換えたとする
openssl dgst -sha256 -verify capub.pem -signature order.sig order.txt    # → Verification Failure
```
署名は秘密鍵を持つ本人しか作れず、検証は公開鍵で誰でもできる。内容が変われば署名と一致しなくなるので、**受信側は改ざんを検知できる**。これがTLS証明書（CAがサーバの公開鍵に署名）や、なりすまし防止の土台になっている。

## 試験での問われ方
- PKIの構成要素（CA・証明書・CRL/OCSP・信頼チェーン）
- サーバ証明書の検証（発行者・有効期限・CN/SAN・失効）、証明書検証不備
- 公開鍵暗号と電子署名、ハッシュ、完全性/真正性/否認防止
- TLSハンドシェイク、HSTS、証明書ピンニング

## 後片付け
```sh
make down LAB=crypto/pki
```
