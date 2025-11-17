#!/bin/bash

# 定义可配置变量
#镜像地址
IMAGE_NAME="hub-dev.cnbn.org.cn/opennet/chainmaker-explorer-backend:v2.3.9"
# 链ID
CHAIN_ID="chainmaker_pk"
# 升级版本
VERSION_RANGE="v2.3.8"

# 执行升级命令
docker run --rm \
  --workdir /chainmaker-explorer-backend \
  --entrypoint "/chainmaker-explorer-backend/bin/upgrader.bin" \
  -v ./cm_explorer_server/config.yml:/chainmaker-explorer-backend/configs/config.yml \
  "$IMAGE_NAME" \
  "$CHAIN_ID" "$VERSION_RANGE" > ./upgrader.log 2>&1