#!/bin/sh
mkdir -p /var/run/vsftpd/empty
/usr/sbin/vsftpd /etc/vsftpd.conf &
/usr/sbin/inetd
nginx
exec sleep infinity
