#!/usr/bin/env bash
#
# Copyright (C) BABEC. All rights reserved.
# Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

function start_container() {
  local container_name=$1
  local image_name=$2
  local port_mapping=$3
  local env_vars=$4

  # 检查容器是否已经存在
  if docker ps -a | grep -q "$container_name"; then 
    echo "$container_name 容器已存在，正在停止并删除..." 
    docker stop "$container_name"
    docker rm "$container_name"
  fi

  # 启动 Docker 容器
  echo "正在启动 $container_name 容器..."
  if docker run --name "$container_name" $env_vars -d -p "$port_mapping" "$image_name"; then
    echo "$container_name 容器启动成功！"
  else
    echo "启动 $container_name 容器失败！"
    exit 1
  fi
}

function start_mysql_container() {
  start_container "ut-mysql-test" "hub-dev.cnbn.org.cn/tools/mysql:8.0.28" "33061:3306" "-e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=chainmaker_explorer_dev"
}

function start_redis_container() {
  start_container "ut-redis-test" "hub-dev.cnbn.org.cn/tools/redis:latest" "63791:6379" ""
}

function ut_cover() {
  # 启动 MySQL Docker 容器
  start_mysql_container
  start_redis_container

  # 等待 docker 启动
  sleep 20
   # 设置环境变量
  export UT_MYSQL_DB_URL="root:123456@tcp(127.0.0.1:33061)/chainmaker_explorer_dev"
  export UT_REDIS_URL="127.0.0.1:63791"
  echo "=====环境变量设置成功！$UT_MYSQL_DB_URL $UT_REDIS_URL===="

  cd ${cm}/$1
  echo "cd ${cm}/$1"
  echo "exec go test"
  go test -coverprofile cover.out ./...
  total=$(go tool cover -func=cover.out | tail -1)
  echo ${total}
  #rm cover.out
  coverage=$(echo ${total} | grep -P '\d+\.\d+(?=\%)' -o) #如果macOS 不支持grep -P选项，可以通过brew install grep更新grep
  #计算注释覆盖率，需要安装gocloc： go install github.com/hhatto/gocloc/cmd/gocloc@latest
  comment_coverage=$(gocloc --include-lang=Go --output-type=json --not-match=".*_test\.go" . | jq '(.total.comment-.total.files*6)/(.total.code+.total.comment)*100')
  echo "注释率：${comment_coverage}%"

  # 如果测试覆盖率低于N，认为ut执行失败
  (( $(awk "BEGIN {print (${coverage} >= $2)}") )) || (echo "$1 单测覆盖率: ${coverage} 低于 $2%"; exit 1)
  (( $(awk "BEGIN {print (${comment_coverage} >= $3)}") )) || (echo "$1 注释覆盖率: ${comment_coverage} 低于 $3%"; exit 1)
}
set -e

cm=$(pwd)

if [[ $cm == *"scripts" ]] ;then
  cm=$cm/..
fi

ut_cover src 50 15

