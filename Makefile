all: proto build docker run-docker

proto:
	sudo docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) cap1573/cap-v3 --proto_path=. --micro_out=. --go_out=:. ./proto/pod/pod.proto

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /usr/local/Cellar/go@1.19/1.19.11/bin/go build -o pod *.go

docker:
	sudo docker build . -t zxnl/pod:latest

run-docker:
	sudo docker run -p 8081:8081 -p 9092:9092 -p 9192:9192 -v /Users/lqy007700/Data/config:/root/.kube/config -v /Users/lqy007700/Data/code/go-application/go-paas/pod/micro.log:/micro.log lqy007700/pod