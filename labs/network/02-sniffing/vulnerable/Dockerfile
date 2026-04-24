# 平文プロトコル（FTP / Telnet / HTTP Basic認証）のサーバ
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      vsftpd nginx openbsd-inetd inetutils-telnetd openssl \
 && rm -rf /var/lib/apt/lists/* /etc/nginx/sites-enabled/default \
 && useradd -m -s /bin/bash alice && echo 'alice:S3cretPass!' | chpasswd \
 && echo 'confidential report' > /home/alice/report.txt \
 && printf 'alice:%s\n' "$(openssl passwd -apr1 'S3cretPass!')" > /etc/nginx/.htpasswd \
 && mkdir -p /var/www/html/secret && echo 'secret page' > /var/www/html/secret/index.html
COPY inetd.conf /etc/inetd.conf
COPY vsftpd.conf /etc/vsftpd.conf
COPY nginx.conf /etc/nginx/conf.d/lab.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
CMD ["/entrypoint.sh"]
