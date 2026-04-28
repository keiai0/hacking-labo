// 脆弱: 境界チェックのない入力でバッファを溢れさせ、隣接する認可フラグを上書きする。
#include <stdio.h>
#include <string.h>

int main(void) {
    // メンバ順を固定するため構造体にする（buf の直後に authorized）
    struct {
        char buf[16];
        int  authorized;
    } s;
    s.authorized = 0;

    printf("パスワードを入力: ");
    scanf("%s", s.buf);                 // 危険: 長さを制限していない

    if (strcmp(s.buf, "correct-horse") == 0)
        s.authorized = 1;

    if (s.authorized)
        puts("ACCESS GRANTED");         // 17文字以上の入力でフラグが上書きされ、ここに到達する
    else
        puts("ACCESS DENIED");
    return 0;
}
