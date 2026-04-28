# C2サーバ役（攻撃者インフラ）。ビーコンを受けて200を返すだけ。
FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends nginx \
 && rm -rf /var/lib/apt/lists/* /etc/nginx/sites-enabled/default
COPY nginx.conf /etc/nginx/conf.d/c2.conf
CMD ["nginx", "-g", "daemon off;"]
