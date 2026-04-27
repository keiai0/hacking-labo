# XXE演習用の小さなPHPサーバ（Goアプリから /xxe としてプロキシされる）
FROM php:8.2-cli-bookworm
COPY xxe.php /var/www/xxe.php
WORKDIR /var/www
EXPOSE 8080
# すべてのリクエストを xxe.php で受ける（ルータスクリプト）
CMD ["php", "-S", "0.0.0.0:8080", "xxe.php"]
