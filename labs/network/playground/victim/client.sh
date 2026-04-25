#!/bin/sh
# 15秒ごとに、平文(web) と 暗号化(secure) の両方へログインする。
USER=alice
PASS='S3cretPass!'
while true; do
  # 平文（web）: FTP / HTTP Basic / Telnet
  curl -s "ftp://web/report.txt" -u "$USER:$PASS" -o /dev/null
  curl -s "http://web/secret/"   -u "$USER:$PASS" -o /dev/null
  /telnet.exp web "$USER" "$PASS" > /dev/null 2>&1
  # 暗号化（secure）: SFTP / HTTPS
  echo 'get report.txt /dev/null' | sshpass -p "$PASS" \
    sftp -q -o BatchMode=no -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -b - "$USER@secure" > /dev/null 2>&1
  curl -sk "https://secure/secret/" -u "$USER:$PASS" -o /dev/null
  sleep 15
done
