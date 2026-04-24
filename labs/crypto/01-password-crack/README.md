# crypto/01-password-crack — パスワード解析（John the Ripper / hashcat）

## 目的
漏洩したハッシュや暗号化ファイルを解析する側を体験し、「保存方式の強さ」で解析コストがどれだけ変わるかを理解する。弱い方式（生MD5・md5crypt・ZipCrypto）は現実的な時間で割れ、強い方式（bcrypt＋強いパスワード）は割れないことを手元で確かめる。

## 構成
ネットワーク攻撃ではなく、入手済みの成果物を解析する想定。生成役が共有ボリュームに成果物を置き、攻撃用コンテナがそれを読む。

| コンテナ | 置き場所（attacker視点） | 中身 |
|---|---|---|
| vuln | `/loot` | `raw_md5.txt`（生MD5）, `shadow`（md5crypt `$1$`）, `secret.zip`（ZipCrypto） |
| fixed | `/loot_fixed` | `shadow`（bcrypt `$2b$`） |

パスワードはいずれも rockyou に含まれる弱いもの（princess など）。fixed の carol だけ強いパスワード。

## 攻撃
```sh
make up    LAB=crypto/01-password-crack
make shell LAB=crypto/01-password-crack
WL=/usr/share/wordlists/rockyou.txt

# 1) 生MD5（ソルトなし）
john --format=raw-md5 --wordlist=$WL /loot/raw_md5.txt
john --format=raw-md5 --show /loot/raw_md5.txt
#   hashcat 版: hashcat -m 0 -a 0 /loot/raw_md5.txt $WL

# 2) shadow（md5crypt）
cp /loot/shadow /tmp/s && john --wordlist=$WL /tmp/s && john --show /tmp/s
#   hashcat 版: hashcat -m 500 -a 0 /loot/shadow $WL

# 3) パスワード付きZIP
zip2john /loot/secret.zip > /tmp/zip.hash
john --wordlist=$WL /tmp/zip.hash && john --show /tmp/zip.hash
```
→ いずれも数秒〜で平文が判明する。

## 対策（fixed との比較）
```sh
cp /loot_fixed/shadow /tmp/sf && john --wordlist=$WL /tmp/sf
```
- bob（弱いパスワードを bcrypt で保存）：割れるが**速度が桁違いに遅い**
- carol（強いパスワードを bcrypt で保存）：**割れない**

同一CPUでの実測（`john --test`／実行時の表示）:

| 方式 | 解析速度の目安 |
|---|---|
| md5crypt（`$1$`） | 約 100,000 回/秒 |
| bcrypt（`$2b$`, cost 12） | 約 40 回/秒 |

約2,500倍遅く、総当たりが現実的でなくなる。さらに「強いパスワード」なら辞書にも載らず割れない。

## 試験での問われ方
- ハッシュ関数とソルトの役割、レインボーテーブル対策
- ストレッチング／計算コスト（bcrypt・scrypt・Argon2・PBKDF2）、パスワードを平文/可逆暗号で保存しない
- 辞書攻撃・総当たり攻撃、パスワードポリシー、アカウントロック
- 「ソルトは秘密ではないが、同一パスワードのハッシュ一致やレインボーテーブルを防ぐ」点

## メモ（M1）
GPUが使えないためCPU解析。弱いパスワード＋小さな前提で十分に成立する。本格的なGPU解析は Mac 上の hashcat を直接使う。

## 後片付け
```sh
make down LAB=crypto/01-password-crack   # 共有ボリューム(loot)も削除される
```
