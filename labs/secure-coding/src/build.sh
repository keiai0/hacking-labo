#!/bin/sh
# 各種コンパイル。/opt/bin に配置。
set -e
mkdir -p /opt/bin
cd /opt/src

# バッファオーバーフロー（保護機能オフ）
gcc -O0 -fno-stack-protector -o /opt/bin/bof_vuln   bof_vuln.c
# 同じ脆弱コードをスタック保護つきでビルド（検知される＝緩和策のデモ）
gcc -O0 -fstack-protector-all -o /opt/bin/bof_canary bof_vuln.c
# 対策版
gcc -O0 -o /opt/bin/bof_fixed bof_fixed.c

# 整数（符号）の検査不備
gcc -O0 -o /opt/bin/intof_vuln  intof_vuln.c
gcc -O0 -o /opt/bin/intof_fixed intof_fixed.c

echo "built: /opt/bin/{bof_vuln,bof_canary,bof_fixed,intof_vuln,intof_fixed}"
