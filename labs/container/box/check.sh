#!/bin/sh
# コンテナの権限・隔離状態を表示する（脱出はしない。設定の良し悪しの観察用）。
echo "== 実行ユーザ =="; id
echo "== 実効ケーパビリティ (CapEff) =="; grep CapEff /proc/self/status
echo "   （人間可読）"; capsh --decode=$(grep CapEff /proc/self/status | awk '{print $2}') 2>/dev/null | sed 's/^/   /'
echo "== ルートFSは書込可能か =="; if touch /_wtest 2>/dev/null; then echo "書込可能(=read_onlyでない)"; rm -f /_wtest; else echo "読み取り専用(read_only)"; fi
echo "== 見えるブロックデバイス（ホストのディスクが見える=危険） =="; ls -l /dev | grep -E '^b' | awk '{print $NF}' | tr '\n' ' '; echo
echo "== /proc/sys は書込可能か（特権だと書ける） =="; (echo 0 > /proc/sys/kernel/dmesg_restrict 2>/dev/null && echo "書込可能(=特権/未制限)") || echo "書込不可(=制限あり)"
echo "== no_new_privs =="; grep NoNewPrivs /proc/self/status
