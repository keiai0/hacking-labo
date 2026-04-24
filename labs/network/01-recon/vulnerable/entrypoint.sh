#!/bin/sh
mkdir -p /run/sshd /var/run/vsftpd/empty
/usr/sbin/sshd
/usr/sbin/vsftpd /etc/vsftpd.conf &
/usr/sbin/inetd
nginx
exec sleep infinity
