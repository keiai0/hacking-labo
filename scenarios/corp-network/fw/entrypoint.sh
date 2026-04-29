#!/bin/sh
set -e
# ip_forward は compose の sysctls でも設定済み。念のため直接も設定（失敗は無視）
echo 1 > /proc/sys/net/ipv4/ip_forward 2>/dev/null || true

EXT=10.10.8.0/24
DMZ=10.10.9.0/24
INT=10.10.10.0/24
WEB=10.10.9.20

iptables -P FORWARD DROP
iptables -A FORWARD -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

# 外部 → DMZ: Web(80) のみ許可。内部へは一切通さない
iptables -A FORWARD -s $EXT -d $WEB -p tcp --dport 80 -j ACCEPT
# DMZ(Web) → 内部: SSH(22)/DB(3306) のみ許可（業務に必要な最小限）
iptables -A FORWARD -s $WEB -d $INT -p tcp -m multiport --dports 22,3306 -j ACCEPT
# それ以外（外部→内部、内部→外部 など）は既定の DROP

echo "FW up: ext→dmz:80, dmz-web→int:22,3306 のみ許可"
exec sleep infinity
