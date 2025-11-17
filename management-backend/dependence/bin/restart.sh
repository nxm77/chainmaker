#
# Copyright (C) BABEC. All rights reserved.
# Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

export LD_LIBRARY_PATH=$(dirname $PWD)/lib:$LD_LIBRARY_PATH
export PATH=$(dirname $PWD)/lib:$PATH
export WASMER_BACKTRACE=1

#support multichain on same machine
#./stop.sh

enable_dockervm={docker_enable}
docker_go_container_name=VM-GO-{org_id}
if [ ${enable_dockervm} == "true" ];then
  which 7z >/dev/null 2>&1
  if [ $? -ne 0 ]; then
  		yum install p7zip p7zip-plugins
  fi
  docker_container_lists=(`docker ps -a | grep ${docker_go_container_name} | awk '{print $1}'`)
  for container_id in ${docker_container_lists[*]}
  do
    docker stop ${container_id}
    docker rm ${container_id}
  done
  ./docker_start.sh
fi

nohup ../../chainmaker start -c ../config/{org_id}/chainmaker.yml > panic.log 2>&1 &
echo "chainmaker is restartting, pls check log..."
