# 対策版: 強い方式（bcrypt、ソルトあり・低速）で保存した成果物を生成する役
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends python3 \
 && rm -rf /var/lib/apt/lists/*
COPY gen.sh /gen.sh
RUN chmod +x /gen.sh
CMD ["/gen.sh"]
