#!/bin/sh
# 弱いパスワード（いずれも rockyou に含まれる）を、弱い方式で保存した成果物を作る。
set -e
OUT=/loot
mkdir -p "$OUT"

# 1) 生のMD5ハッシュ（ソルトなし）: hashcat -m 0 / john raw-md5
: > "$OUT/raw_md5.txt"
for p in password123 iloveyou sunshine princess; do
  printf '%s' "$p" | python3 -c 'import sys,hashlib;print(hashlib.md5(sys.stdin.buffer.read()).hexdigest())' >> "$OUT/raw_md5.txt"
done

# 2) shadow ファイル（md5crypt $1$、ソルトはあるが高速で危険）: john / hashcat -m 500
python3 - <<'PY' > "$OUT/shadow"
import crypt
users = {"bob":"iloveyou", "carol":"sunshine", "dave":"princess"}
for u,p in users.items():
    print(f"{u}:{crypt.crypt(p, crypt.mksalt(crypt.METHOD_MD5))}:19000:0:99999:7:::")
PY

# 3) パスワード付きZIP（従来方式ZipCrypto、弱いパスワード）: zip2john → john
echo 'TOP SECRET: launch code 0000' > /tmp/secret.txt
( cd /tmp && zip -q -P sunshine "$OUT/secret.zip" secret.txt )

chmod -R a+r "$OUT"
echo "[vulnerable] 生成完了: raw_md5.txt / shadow / secret.zip"
exec sleep infinity
