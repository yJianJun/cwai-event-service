.PHONY: build run import

build:
	go env -w GO111MODULE="on"
	go env -w GOPROXY="https://goproxy.cn,direct"
	go env -w GOPRIVATE=""
	go build -o server cmd/server/server.go

run:
	./server -c conf/debug.yaml

import:
	goimports -local server-template/ -d -w ./
