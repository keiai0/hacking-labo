# auth — OAuth 2.0 / OIDC の設定・実装ミス（軽量IdP+RP）

Keycloakを使わず、自作の小さな認可サーバ(IdP)＋RPで、OAuth/OIDCの典型的な設定・実装ミスを再現する。標準ライブラリのみ・即起動。各エンドポイントは `?mode=vuln|fixed`。

## 構成
- `app`（127.0.0.1:8083）: IdP（/authorize, /token）＋ RP（/rp/login, /rp/callback, /rp/validate）＋ 攻撃者受け皿（/attacker/catch）
- 登録済み: `client_id=corp-portal` / `redirect_uri=http://client.corp/callback`

```sh
make up LAB=auth
# ブラウザ: http://127.0.0.1:8083/ （Burpでフローを観察）
```

## 1) redirect_uri 検証不備（認可コード横取り）
```sh
B=http://127.0.0.1:8083
# 脆弱: 前方一致 → client.corp.evil が通り、コードが攻撃者へ
curl -si "$B/authorize?mode=vuln&client_id=corp-portal&redirect_uri=http://client.corp.evil/cb&response_type=code" | grep -i location
# 対策: 完全一致 → 400
curl -s -o /dev/null -w '%{http_code}\n' "$B/authorize?mode=fixed&client_id=corp-portal&redirect_uri=http://client.corp.evil/cb&response_type=code"
```
対策: redirect_uri は登録値と**完全一致**で検証（部分/前方一致は不可）。

## 2) state 欠如（CSRF / コード注入）
```sh
# 脆弱: stateを検証しない → 攻撃者が用意したcodeで被害者をログインさせられる
curl -s "$B/rp/callback?mode=vuln&code=ATTACKERCODE"
# 対策: state不一致は拒否（400）
curl -s -o /dev/null -w '%{http_code}\n' "$B/rp/callback?mode=fixed&code=ATTACKERCODE"
```
対策: 認可リクエストで `state` を発行し、コールバックで一致を検証（`/rp/login?mode=fixed` はstate/PKCEを付与）。

## 3) PKCE なし（横取りコードの悪用）
```sh
V=verifier-xyz
CHAL=$(printf '%s' "$V" | openssl dgst -sha256 -binary | openssl base64 -A | tr '+/' '-_' | tr -d '=')
getcode(){ curl -s -o /dev/null -w '%{redirect_url}' "$B/authorize?mode=$1&client_id=corp-portal&redirect_uri=http://client.corp/callback&response_type=code&code_challenge=$CHAL" | sed 's/.*code=//;s/&.*//'; }
# 脆弱: code_verifier を検証しない → 横取りcodeをそのまま交換できる
curl -s -o /dev/null -w 'vuln=%{http_code}\n' "$B/token?mode=vuln&code=$(getcode vuln)&redirect_uri=http://client.corp/callback"
# 対策: verifier無しは失敗(400)、正しいverifierなら成功(200)
curl -s -o /dev/null -w 'fixed(無)=%{http_code}\n' "$B/token?mode=fixed&code=$(getcode fixed)&redirect_uri=http://client.corp/callback"
curl -s -o /dev/null -w 'fixed(正)=%{http_code}\n' "$B/token?mode=fixed&code=$(getcode fixed)&redirect_uri=http://client.corp/callback&code_verifier=$V"
```
対策: PKCE（S256）を必須化し、`code_verifier` を検証。

## 4) IDトークン検証不備
```sh
b64u(){ printf '%s' "$1" | openssl base64 -A | tr '+/' '-_' | tr -d '='; }
FORGED="$(b64u '{"alg":"none","typ":"JWT"}').$(b64u '{"iss":"http://authz","aud":"corp-portal","sub":"attacker","role":"admin"}')."
curl -s   "$B/rp/validate?mode=vuln&token=$FORGED"   # 脆弱: alg=none/無検証で admin を受理
curl -s -o /dev/null -w '%{http_code}\n' "$B/rp/validate?mode=fixed&token=$FORGED"  # 対策: 401
```
対策: 署名（alg固定・鍵で検証）＋ `iss`/`aud`/有効期限 を検証。`alg=none` を拒否。

## 試験での問われ方
- 認可コードフロー、`state`（CSRF対策）、PKCE、`redirect_uri` 厳密一致
- アクセストークンとIDトークンの違い、IDトークン(JWT)の検証項目（署名/iss/aud/exp/nonce）
- 暗黙フローの非推奨、スコープ最小化、トークンの保管
- OAuth(認可) と OIDC(認証) の役割分担

## 補足
本ラボは要点を自作で最小再現したもの。実製品の設定項目（Keycloakのクライアント設定等）を触りたい場合は、別途Keycloakを立てて確認する。

## 後片付け
```sh
make down LAB=auth
```
