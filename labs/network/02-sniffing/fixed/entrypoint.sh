#!/bin/sh
mkdir -p /run/sshd
/usr/sbin/sshd
nginx
exec sleep infinity
