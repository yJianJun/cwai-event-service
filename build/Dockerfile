FROM golang:1.20.2 as base

WORKDIR /go/src/ctyun-code.srdcloud.cn/aiplat/cwai-watcher
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o  cwai-watcher  cmd/server/server.go

FROM alpine:3.12

COPY --from=base /go/src/ctyun-code.srdcloud.cn/aiplat/cwai-watcher/cwai-watcher /bin/

WORKDIR /bin
CMD ["/bin/cwai-watcher","-c","/etc/cwai-watcher/config.yaml"]
