#!/bin/sh
# GW 経由で社内の他セグメントへ経路を追加（Dockerの既定GWとは別のFWを使う）
# ROUTES="10.10.9.0/24 10.10.10.0/24" GW=10.10.8.2 のように指定
if [ -n "$GW" ] && [ -n "$ROUTES" ]; then
  for net in $ROUTES; do
    ip route replace "$net" via "$GW" 2>/dev/null || true
  done
fi

case " $START " in
  *" sshd "*) mkdir -p /run/sshd; /usr/sbin/sshd ;;
esac
case " $START " in
  *" nginx "*)
    echo "<h1>${HOSTNAME} (DMZ Web)</h1>" > /var/www/html/index.html 2>/dev/null || true
    nginx ;;
esac

# 内部サーバには「機密ファイル」を置く（持ち出しの的）
if [ -n "$SECRET" ]; then echo "$SECRET" > /root/crown-jewels.txt; fi

exec sleep infinity
