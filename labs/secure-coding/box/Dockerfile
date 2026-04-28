# C言語のセキュアコーディング演習（gcc/gdb）。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      gcc libc6-dev gdb make vim-tiny \
 && rm -rf /var/lib/apt/lists/*
COPY src/ /opt/src/
RUN chmod +x /opt/src/build.sh && /opt/src/build.sh
WORKDIR /opt
CMD ["sleep", "infinity"]
