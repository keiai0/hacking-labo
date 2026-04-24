# 暗号化プロトコル（SSH/SFTP / HTTPS）のサーバ
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      openssh-server nginx openssl \
 && rm -rf /var/lib/apt/lists/* /etc/nginx/sites-enabled/default \
 && useradd -m -s /bin/bash alice && echo 'alice:S3cretPass!' | chpasswd \
 && echo 'confidential report' > /home/alice/report.txt \
 && printf 'alice:%s\n' "$(openssl passwd -apr1 'S3cretPass!')" > /etc/nginx/.htpasswd \
 && mkdir -p /var/www/html/secret && echo 'secret page' > /var/www/html/secret/index.html \
 && openssl req -x509 -nodes -newkey rsa:2048 -days 3650 -subj /CN=fixed \
      -keyout /etc/nginx/server.key -out /etc/nginx/server.crt
COPY nginx.conf /etc/nginx/conf.d/lab.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
CMD ["/entrypoint.sh"]
