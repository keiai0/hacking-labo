#!/bin/sh
# 認証ログを生成。正常ログインに混じって、特定IPからの総当たり（多数のFailed）を出す。
mkdir -p /logs
log() { echo "$(date '+%b %d %H:%M:%S') corp sshd: $1" >> /logs/auth.log; }
while true; do
  log "Accepted password for alice from 10.0.0.5 port 51515 ssh2"
  n=0
  while [ $n -lt 8 ]; do                       # 攻撃者IPからの連続失敗（総当たりの兆候）
    log "Failed password for admin from 203.0.113.9 port 40000 ssh2"
    n=$((n + 1))
  done
  log "Accepted password for bob from 10.0.0.6 port 51516 ssh2"
  sleep 5
done
