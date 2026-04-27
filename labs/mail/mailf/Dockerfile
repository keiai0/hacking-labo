# 対策版メールサーバ: 中継制限 + SPF検証(policyd-spf)
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends postfix postfix-policyd-spf-python \
 && rm -rf /var/lib/apt/lists/* \
 && useradd -m alice \
 && (id policyd-spf 2>/dev/null || useradd -r -s /usr/sbin/nologin policyd-spf)
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
CMD ["/entrypoint.sh"]
