#!/bin/sh
set -e
postconf -e \
  "myhostname=mail.corp.local" \
  "mydomain=corp.local" \
  "mydestination=corp.local, localhost" \
  "inet_interfaces=all" \
  "inet_protocols=ipv4" \
  "mynetworks=127.0.0.0/8" \
  "smtpd_recipient_restrictions=permit_mynetworks, reject_unauth_destination, check_policy_service unix:private/policyd-spf" \
  "policyd-spf_time_limit=3600" \
  "maillog_file=/dev/stdout" \
  "compatibility_level=3.6"
# reject_unauth_destination → 中継を自ドメイン宛に限定（オープンリレーでない）
# check_policy_service policyd-spf → 送信元ドメインのSPFを検証し、失敗を拒否
postconf -M "policyd-spf/unix=policyd-spf unix - n n - 0 spawn user=policyd-spf argv=/usr/bin/policyd-spf"
exec postfix start-fg
