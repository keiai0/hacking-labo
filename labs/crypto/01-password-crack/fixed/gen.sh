#!/bin/sh
# 対策版: 同じ利用者を bcrypt（ソルトあり・計算コスト高）で保存する。
# 弱いパスワードでも1件ずつが遅く、強いパスワードは現実的に割れない。
set -e
OUT=/loot_fixed
mkdir -p "$OUT"

python3 - <<'PY' > "$OUT/shadow"
import crypt
# bob は弱いパスワードのまま（方式だけ強く）/ carol は強いパスワード
users = {"bob":"iloveyou", "carol":"S7!qZ-watermelon-42"}
for u,p in users.items():
    h = crypt.crypt(p, crypt.mksalt(crypt.METHOD_BLOWFISH))  # $2b$ = bcrypt
    print(f"{u}:{h}:19000:0:99999:7:::")
PY

chmod -R a+r "$OUT"
echo "[fixed] 生成完了: shadow（bcrypt）"
exec sleep infinity
