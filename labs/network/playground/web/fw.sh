#!/bin/sh
# ファイアウォールのトグル。make fw-on / fw-off から呼ばれる。
case "$1" in
  on)
    iptables -P INPUT DROP
    iptables -A INPUT -i lo -j ACCEPT
    iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
    iptables -A INPUT -p tcp -m multiport --dports 22,80 -j ACCEPT
    echo "FW ON: 22,80 のみ許可（他はDROP）"
    ;;
  off)
    iptables -F; iptables -P INPUT ACCEPT
    echo "FW OFF: すべて許可"
    ;;
  *) echo "usage: fw.sh on|off"; exit 1 ;;
esac
