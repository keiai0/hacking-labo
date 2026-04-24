# 攻撃用コンテナ（全ラボ共通）
FROM kalilinux/kali-rolling@sha256:7997e0581c4f5c185cfe6773d5cae8fcc7c6b79798726ccaabf843af4bcd2013

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
      # 基本
      bash-completion ca-certificates curl wget vim-tiny less jq procps \
      iproute2 iputils-ping net-tools dnsutils netcat-openbsd arping \
      telnet ftp openssh-client sshpass \
      # 偵察
      nmap \
      # 盗聴・MITM
      tcpdump tshark dsniff mitmproxy python3-scapy \
      # パスワード
      hydra john hashcat wordlists \
      # Web
      sqlmap ffuf nikto \
      # 暗号・メール
      openssl swaks \
      # スクリプト
      python3 python3-requests python3-pwntools \
 && gzip -d /usr/share/wordlists/rockyou.txt.gz \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /work
CMD ["sleep", "infinity"]
