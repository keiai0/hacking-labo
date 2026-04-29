#!/bin/sh
# 自前CA・正規サーバ証明書・なりすまし証明書を生成して /pki に出力する。
set -e
cd /pki
if [ -f ca.crt ] && [ -f legit.crt ] && [ -f mitm.crt ]; then
  echo "already generated"; exit 0
fi

# 1) 自前の認証局(CA)
openssl req -x509 -newkey rsa:2048 -nodes -days 3650 \
  -subj "/CN=HackLabo Root CA" -keyout ca.key -out ca.crt

# 2) 正規サーバ(bank.local)の証明書：CAで署名
printf "subjectAltName=DNS:bank.local\n" > /tmp/san.ext
openssl req -newkey rsa:2048 -nodes -subj "/CN=bank.local" -keyout legit.key -out legit.csr
openssl x509 -req -in legit.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -days 3650 -extfile /tmp/san.ext -out legit.crt

# 3) なりすましサーバ：同じCN=bank.localだが「自己署名」（CAの署名なし）
openssl req -x509 -newkey rsa:2048 -nodes -days 3650 \
  -subj "/CN=bank.local" -addext "subjectAltName=DNS:bank.local" \
  -keyout mitm.key -out mitm.crt

chmod a+r /pki/*.key /pki/*.crt
echo "generated: ca / legit / mitm"
