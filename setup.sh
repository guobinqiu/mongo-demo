#!/bin/bash
set -e # 脚本遇到错误会立即停止，避免继续执行

echo "停止并删除旧容器..."
docker rm -f mongo >/dev/null 2>&1 || true

echo "启动 MongoDB 副本集容器..."
docker run -d \
  --name mongo \
  -p 27017:27017 \
  mongo:8 \
  --replSet rs0

echo "等待MongoDB启动完成..."
sleep 5

echo "初始化副本集（host: localhost:27017）..."
docker exec -it mongo mongosh --eval '
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "localhost:27017" }
  ]
})
'

echo "MongoDB 副本集已启动并配置完成！"
echo "连接字符串请使用：mongodb://localhost:27017/?replicaSet=rs0/?replicaSet=rs0"
