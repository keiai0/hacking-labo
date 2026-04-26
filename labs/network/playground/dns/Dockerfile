# DNSサーバ役（BIND9）。ゾーン転送(AXFR)が誰でも可能な設定ミスを題材にする。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends bind9 \
 && rm -rf /var/lib/apt/lists/*
COPY named.conf /etc/bind/named.conf
COPY db.corp.local /etc/bind/db.corp.local
CMD ["/usr/sbin/named", "-g", "-c", "/etc/bind/named.conf"]
