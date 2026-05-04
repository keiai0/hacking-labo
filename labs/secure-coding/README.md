# secure-coding — C言語のセキュアコーディング（BOF・整数検査）

午後試験のコード読解対策。境界外書き込みと整数（符号）の検査不備を、脆弱版/対策版のCコードで比較する。

## 構成
`box` コンテナ（gcc/gdb）で、ビルド済みバイナリ `/opt/bin/*` を動かす。ネットワーク不要。
```sh
make up    LAB=secure-coding
make shell LAB=secure-coding SVC=box      # コンテナに入る
```
ソースは `src/`（ホストからも編集可、再ビルドは `/opt/src/build.sh`）。

## 1) バッファオーバーフロー（認証バイパス）
**目的**: 長さを制限しない入力関数が、確保したバッファを越えて隣接する変数まで書き換えてしまい、攻撃者がパスワードを知らなくても認可フラグを立てて認証を回避できることを、脆弱版と対策版で比較する。

`bof_vuln.c`: `char buf[16]` の直後に `int authorized`。`scanf("%s")` が長さを制限しないため、17文字以上で（buf を溢れさせて）隣の `authorized` を上書きできる。
```sh
echo correct-horse        | /opt/bin/bof_vuln   # 正規 → ACCESS GRANTED
echo AAAAAAAAAAAAAAAAAAAA | /opt/bin/bof_vuln   # 20文字 → ACCESS GRANTED（パスワード不要でバイパス）
echo AAAAAAAAAAAAAAAAAAAA | /opt/bin/bof_fixed  # 対策(fgets) → ACCESS DENIED
```
対策: `fgets(buf, sizeof(buf), stdin)` で長さを制限。

### 緩和策: スタックカナリア
同じ脆弱コードを `-fstack-protector-all` でビルドした `bof_canary`:
```sh
head -c 64 /dev/zero | tr '\0' A | /opt/bin/bof_canary   # *** stack smashing detected *** → abort
```
※ カナリアは「戻りアドレスの破壊」を復帰時に検知する。隣接ローカル変数の上書き（上の20文字）は検知しない点に注意（コンパイラ緩和は多層防御の一部で、根本対策はコードの境界チェック）。

## 2) 整数（符号）の検査不備
**目的**: 符号付き整数で長さを受けて上限だけ検査すると、攻撃者が負数を入れて検査をすり抜け、符号なし(`size_t`)への変換で巨大値に化けて境界外アクセス（クラッシュ/メモリ破壊）を起こせることを確認する。

`intof_vuln.c`: 長さを `int` で受け、上限だけ検査。負数が上限チェックを通り、`memcpy` で `size_t` の巨大値になって境界外アクセス。
```sh
printf '%s\n' -1 | /opt/bin/intof_vuln   # Segmentation fault（検査をすり抜けクラッシュ）
printf '%s\n' -1 | /opt/bin/intof_fixed  # 不正なサイズです（下限も検査）
printf '%s\n' 32 | /opt/bin/intof_vuln   # 正常
```
対策: 符号・下限・上限をすべて検査（`n >= 0 && n <= sizeof(dst)`）。

## gdbで観察（任意）
```sh
gdb -q /opt/bin/bof_vuln      # run、スタックやフレームを確認
```
※ M1（arm64）なのでレジスタ名等は教科書のx86と異なるが、「境界を越えて隣の領域を壊す」原理は同じ。

## 試験での問われ方
- バッファオーバーフロー・境界チェック、安全な関数（`fgets`/`strncpy`/`snprintf`）
- 整数オーバーフロー・符号の扱い・サイズ計算の検査
- 緩和策（スタックカナリア / DEP(NX) / ASLR / RELRO）と、その限界
- 競合状態（TOCTOU）※別途追加予定

## 後片付け
```sh
make down LAB=secure-coding
```
