#!/bin/sh
# 対策1: 必要なサービス（SSH, HTTP）だけを起動
mkdir -p /run/sshd
/usr/sbin/sshd
nginx

# 対策2: ホストFWで許可ポート以外を破棄（DROP）
iptables -P INPUT DROP
iptables -A INPUT -i lo -j ACCEPT
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -p tcp -m multiport --dports 22,80 -j ACCEPT

exec sleep infinity
