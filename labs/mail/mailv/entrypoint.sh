#!/bin/sh
set -e
postconf -e \
  "myhostname=mail.corp.local" \
  "mydomain=corp.local" \
  "mydestination=corp.local, localhost" \
  "inet_interfaces=all" \
  "inet_protocols=ipv4" \
  "mynetworks=127.0.0.0/8" \
  "smtpd_recipient_restrictions=permit" \
  "maillog_file=/dev/stdout" \
  "compatibility_level=3.6"
# smtpd_recipient_restrictions=permit → 誰からでも任意宛先を中継（オープンリレー）
# SPF/DMARCの検証は一切しない
exec postfix start-fg
