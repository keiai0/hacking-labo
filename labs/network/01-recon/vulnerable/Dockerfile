FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
RUN apt-get update && apt-get install -y --no-install-recommends \
      openssh-server vsftpd nginx openbsd-inetd inetutils-telnetd \
 && rm -rf /var/lib/apt/lists/* /etc/nginx/sites-enabled/default \
 && useradd -m -s /bin/bash alice && echo 'alice:password123' | chpasswd
COPY inetd.conf /etc/inetd.conf
COPY vsftpd.conf /etc/vsftpd.conf
COPY nginx.conf /etc/nginx/conf.d/lab.conf
COPY www/ /var/www/
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
CMD ["/entrypoint.sh"]
