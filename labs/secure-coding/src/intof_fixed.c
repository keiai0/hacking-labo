// 対策: 符号なしで受け取り、0以上かつバッファ内であることを検査する。
#include <stdio.h>
#include <string.h>

int main(void) {
    char src[256];
    char dst[64];
    memset(src, 'A', sizeof(src));

    long n;
    printf("コピーするバイト数: ");
    if (scanf("%ld", &n) != 1) return 1;

    if (n >= 0 && n <= (long)sizeof(dst)) {   // 下限と上限の両方を検査
        memcpy(dst, src, (size_t)n);
        printf("%ld バイトをコピーしました\n", n);
    } else {
        puts("不正なサイズです");
    }
    return 0;
}
