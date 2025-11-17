GIT_BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)

build:
	go mod tidy && cd src && go build -o ../scripts/chainmaker-browser.bin && cd upgrader && go build -o ../../scripts/upgrader.bin

run:
	go mod tidy && cd src  && go build -o ../scripts/chainmaker-browser.bin  && cd ../ && ./scripts/chainmaker-browser.bin > web.log 2>&1

run-dev:
	go mod tidy && cd src  && go build -o ../scripts/chainmaker-browser.bin  && cd ../ && ./scripts/chainmaker-browser.bin --env dev > web.log 2>&1

run-cert:
	go mod tidy && cd src  && go build -o ../scripts/chainmaker-browser.bin  && cd ../ && ./scripts/chainmaker-browser.bin --env cert

run-upgrader:
	go mod tidy && cd src  && cd upgrader && go build -o ../../scripts/upgrader.bin  && cd ../../ && ./scripts/upgrader.bin chain1 v2.3.8 > upgrader.log 2>&1


docker-stop-mysql:
	cd docker && docker-compose -f docker-compose.yml down

docker-start-mysql:
	cd docker && docker-compose -f docker-compose.yml up -d

docker-stop-clickhouse-v:
	cd docker && docker-compose -f docker-compose-clickhouse.yml down -v

docker-stop-clickhouse:
	cd docker && docker-compose -f docker-compose-clickhouse.yml down

docker-start-clickhouse:
	cd docker && docker-compose -f docker-compose-clickhouse.yml up -d

docker-build:
	docker build . -t chainmaker-explorer-backend:latest -f Dockerfile

docker-push:
	docker image tag chainmaker-explorer-backend hub-dev.cnbn.org.cn/opennet/chainmaker-explorer-backend:${GIT_BRACH_NAME}
	docker push hub-dev.cnbn.org.cn/opennet/chainmaker-explorer-backend:${GIT_BRACH_NAME}
	docker image tag chainmaker-explorer-backend chainmaker1.tencentcloudcr.com/opennet/chainmaker-explorer-backend:${GIT_BRACH_NAME}
	docker push chainmaker1.tencentcloudcr.com/opennet/chainmaker-explorer-backend:${GIT_BRACH_NAME}

ut:
	./scripts/ut_cover.sh

lint:
	golangci-lint run ./...

docker_kub_pod:
	kubectl get pod  -o wide | grep  explorer-backend

#算力平台主链configmap
k8s_configmap:
	kubectl create cm explorer-backend-client  -n default --from-file=./k8s/configs/crypto-config/node1/user/client1/client1.key --dry-run=client -o yaml > ./k8s/explorer-backend-client.yaml
	kubectl create cm explorer-backend-config  -n default --from-file=./k8s/configs/config.yml --dry-run=client -o yaml > ./k8s/explorer-backend-config.yaml
	kubectl apply -f ./k8s/explorer-backend-client.yaml
	kubectl apply -f ./k8s/explorer-backend-config.yaml

#算力平台主链
k8s_apply:
	kubectl delete -f ./k8s/deploy-explorer-backend.yaml
	kubectl apply -f ./k8s/deploy-explorer-backend.yaml

#k8s重启前端服务
k8s_restart_front_formal:
	kubectl rollout restart deployment/explorer-front-formal

zip:
	cd k8s && zip -r ../explorer_config.zip ./configs && cd ../
