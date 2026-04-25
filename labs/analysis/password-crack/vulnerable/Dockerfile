# 「漏洩した」弱いハッシュ・成果物を生成して共有ボリュームに置く役
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      python3 zip coreutils \
 && rm -rf /var/lib/apt/lists/*
COPY gen.sh /gen.sh
RUN chmod +x /gen.sh
CMD ["/gen.sh"]
