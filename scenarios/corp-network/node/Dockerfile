# 汎用ノード（DMZ/内部のサーバ役）。env で経路とサービスを切り替える。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      iproute2 iputils-ping curl netcat-openbsd openssh-server nginx \
 && rm -rf /var/lib/apt/lists/* /etc/nginx/sites-enabled/default \
 && useradd -m -s /bin/bash svc && echo 'svc:Passw0rd!' | chpasswd \
 && mkdir -p /var/www/html
COPY web.conf /etc/nginx/conf.d/web.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
CMD ["/entrypoint.sh"]
