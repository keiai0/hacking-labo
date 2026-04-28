// 脆弱: 符号付きintで長さを検査。負数を渡すと上限チェックをすり抜け、
// memcpy では巨大な size_t に変換されて境界外アクセス（クラッシュ/破壊）。
#include <stdio.h>
#include <string.h>

int main(void) {
    char src[256];
    char dst[64];
    memset(src, 'A', sizeof(src));

    int n;
    printf("コピーするバイト数: ");
    if (scanf("%d", &n) != 1) return 1;

    if (n <= (int)sizeof(dst)) {        // 危険: n<0 もここを通る
        memcpy(dst, src, n);            // n=-1 → size_tで巨大値に
        printf("%d バイトをコピーしました\n", n);
    } else {
        puts("大きすぎます");
    }
    return 0;
}
