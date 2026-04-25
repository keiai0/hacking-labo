# 利用者役: 定期的に各サーバへログインし、盗聴・MITMの題材になる通信を流す。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      curl expect telnet openssh-client sshpass ca-certificates \
 && rm -rf /var/lib/apt/lists/*
COPY client.sh telnet.exp /
RUN chmod +x /client.sh /telnet.exp
CMD ["/client.sh"]
