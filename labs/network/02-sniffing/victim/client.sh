#!/bin/sh
# 15秒ごとに平文 / 暗号化の両方でログインする
USER=alice
PASS='S3cretPass!'

while true; do
  # 平文
  curl -s "ftp://vuln/report.txt" -u "$USER:$PASS" -o /dev/null
  curl -s "http://vuln/secret/"   -u "$USER:$PASS" -o /dev/null
  /telnet.exp vuln "$USER" "$PASS" > /dev/null 2>&1

  # 暗号化
  echo 'get report.txt /dev/null' | sshpass -p "$PASS" \
    sftp -q -o BatchMode=no -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -b - "$USER@fixed" > /dev/null 2>&1
  curl -sk "https://fixed/secret/" -u "$USER:$PASS" -o /dev/null

  sleep 15
done
