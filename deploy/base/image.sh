# -n for nano
docker buildx build --platform linux/amd64,linux/arm64 -t harbor.newcapec.cn/cncamp/nfs-ganesha:v3.5-n ./ --push
