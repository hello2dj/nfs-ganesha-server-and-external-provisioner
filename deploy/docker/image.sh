#!/bin/sh
docker buildx build --platform linux/amd64,linux/arm64 -t harbor.newcapec.cn/cncamp/qy-nfs-server:v3.5-7 ./ --push
