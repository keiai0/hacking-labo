# コンテナ設定の良し悪しを観察するための小さなイメージ。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      libcap2-bin procps mount \
 && rm -rf /var/lib/apt/lists/* \
 && useradd -u 1000 -m appuser
COPY check.sh /check.sh
RUN chmod +x /check.sh
CMD ["sleep", "infinity"]
