# container — コンテナのセキュリティ（特権/権限）とクラウド設定ミス

コンテナの設定ミス（特権付与・ケーパビリティ過多）がどれだけ攻撃面を広げるかを**観察**で理解する。実際にホストへ脱出する手順は、ホストを汚さないため**使い捨てVM内**で行う（下記）。

## A. 設定ミス vs ハードニングの比較（安全・ホスト不可侵）
```sh
make up LAB=container
docker compose -f labs/container/compose.yml -p hl-container exec misconfig /check.sh
docker compose -f labs/container/compose.yml -p hl-container exec hardened  /check.sh
```

| 項目 | misconfig（`privileged: true`） | hardened |
|---|---|---|
| 実行ユーザ | root | 非root（uid 1000） |
| ケーパビリティ | 全付与（cap_sys_admin, cap_sys_module …） | 全剥奪（CapEff=0） |
| ルートFS | 書込可能 | 読み取り専用 |
| ブロックデバイス | **ホストのディスク(vda等)が見える** | 見えない |
| /proc/sys | 書込可能 | 不可 |
| no_new_privs | 無効 | 有効 |

特権コンテナから**ホストのディスクが見える**ことが、脱出（そのディスクをマウントして読む等）の起点になる。ハードニング（`cap_drop: ALL` / `read_only` / `no-new-privileges` / 非root / `privileged`禁止）で攻撃面が消える。

### ハードニングの指針（compose.yml 参照）
- `privileged` を使わない。必要な権限だけ `cap_add`
- `cap_drop: ["ALL"]`、`security_opt: ["no-new-privileges:true"]`
- `read_only: true` ＋ 必要箇所のみ `tmpfs`
- `user:` で非rootに、`docker.sock` をマウントしない

## B. 実際の脱出（使い捨てVM内でのみ）
ホスト（Mac/OrbStack）を汚さないよう、専用VMで行う。
```sh
orb create ubuntu escape-lab        # 使い捨てVM
orb -m escape-lab                   # VMに入る（この中にDockerを入れて実施）
```
VM内で試す代表的な脱出:
1. **特権＋ホストディスクのマウント**: `privileged` コンテナ内で `mount /dev/vda1 /mnt` → ホストのファイルを読み書き。
2. **docker.sock のマウント**: `/var/run/docker.sock` を持つコンテナから `docker run -v /:/host ...` でホストFSを掌握。
3. **`cap_sys_admin` + cgroup release_agent** などの既知手法。

※ これらは本物のホストに対しては実行しない。VMは終わったら `orb delete escape-lab` で破棄。

## C. クラウド設定ミス（任意・軽量）
- **公開ストレージ**: MinIO を立て、バケットに匿名read可能ポリシーを付与 → 認証なしで取得できることを確認（LocalStack/MinIO、後日追加）。
- **メタデータSSRF**: `169.254.169.254` 風のモックを立て、webappのSSRFから一時認証情報を抜く筋書き（webapp と連携）。

## 試験での問われ方
- コンテナの隔離（名前空間・cgroup・ケーパビリティ）と最小権限
- 特権コンテナ/ホストマウント/docker.sock の危険性
- イメージの脆弱性管理、ルートレス、読み取り専用FS
- クラウドの共有責任モデル、公開設定・IAM過剰権限・メタデータ悪用

## 後片付け
```sh
make down LAB=container
```
