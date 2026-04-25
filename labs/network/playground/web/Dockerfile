# 社内サーバ役: 平文サービス（FTP/Telnet/HTTP/HTTP Basic）+ 管理画面 + SSH
# recon（偵察）と sniffing（盗聴）の標的を兼ねる。FWトグルで防御の有無も切替。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      openssh-server vsftpd nginx openbsd-inetd inetutils-telnetd openssl iptables \
 && rm -rf /var/lib/apt/lists/* /etc/nginx/sites-enabled/default \
 && useradd -m -s /bin/bash alice && echo 'alice:S3cretPass!' | chpasswd \
 && echo 'confidential report' > /home/alice/report.txt \
 && printf 'alice:%s\n' "$(openssl passwd -apr1 'S3cretPass!')" > /etc/nginx/.htpasswd \
 && mkdir -p /var/www/html/backup /var/www/html/secret /var/www/admin \
 && echo '<h1>Corp Web</h1>' > /var/www/html/index.html \
 && printf -- "-- users dump\nINSERT INTO users VALUES ('alice','S3cretPass!');\n" > /var/www/html/backup/db_2026.sql \
 && echo 'secret page' > /var/www/html/secret/index.html \
 && echo '<h1>Admin Console</h1><p>認証なし</p>' > /var/www/admin/index.html
COPY inetd.conf /etc/inetd.conf
COPY vsftpd.conf /etc/vsftpd.conf
COPY nginx.conf /etc/nginx/conf.d/lab.conf
COPY fw.sh /fw.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /fw.sh /entrypoint.sh
CMD ["/entrypoint.sh"]
