FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      openssh-server nginx iptables \
 && rm -rf /var/lib/apt/lists/* /etc/nginx/sites-enabled/default \
 && useradd -m -s /bin/bash alice && echo 'alice:password123' | chpasswd \
 && echo 'DebianBanner no' >> /etc/ssh/sshd_config \
 && mkdir -p /var/www/html && echo '<h1>Corp Web</h1>' > /var/www/html/index.html
COPY nginx.conf /etc/nginx/conf.d/lab.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
CMD ["/entrypoint.sh"]
