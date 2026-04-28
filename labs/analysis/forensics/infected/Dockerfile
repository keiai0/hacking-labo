# 感染ホスト役。C2へ一定間隔でビーコン（情報を持ち出す模擬通信）を送る。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends curl \
 && rm -rf /var/lib/apt/lists/*
COPY beacon.sh /beacon.sh
RUN chmod +x /beacon.sh
CMD ["/beacon.sh"]
