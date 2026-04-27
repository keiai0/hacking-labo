# メール演習用DNS。partner.example の SPF/DMARC レコードを提供する。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends bind9 \
 && rm -rf /var/lib/apt/lists/*
COPY named.conf /etc/bind/named.conf
COPY db.partner.example /etc/bind/db.partner.example
CMD ["/usr/sbin/named", "-g", "-c", "/etc/bind/named.conf"]
