mkdir -p redis-config;

for i in `seq 1 6`; do
    cp "cluster-conf.conf" "redis-config/cluster-conf${i}.conf";
    sed -i "s/\${PORT}/638${i}/g" "redis-config/cluster-conf${i}.conf";
    echo "生成文件:redis-config/cluster-conf${i}.conf";
done;

docker compose up -d;

# 建立集群
# /usr/local/bin/redis-cli --cluster create \
# 172.30.28.83:6381 172.30.28.83:6382 172.30.28.83:6383 \
# 172.30.28.83:6384 172.30.28.83:6385 172.30.28.83:6386 \
# --cluster-replicas 1;