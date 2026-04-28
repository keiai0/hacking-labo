// 対策: 長さを制限して読み込む（fgets + バッファサイズ）。
#include <stdio.h>
#include <string.h>

int main(void) {
    struct {
        char buf[16];
        int  authorized;
    } s;
    s.authorized = 0;

    printf("パスワードを入力: ");
    if (fgets(s.buf, sizeof(s.buf), stdin) == NULL) return 1;  // サイズ超過を切り詰める
    s.buf[strcspn(s.buf, "\n")] = '\0';                        // 改行除去

    if (strcmp(s.buf, "correct-horse") == 0)
        s.authorized = 1;

    puts(s.authorized ? "ACCESS GRANTED" : "ACCESS DENIED");
    return 0;
}
