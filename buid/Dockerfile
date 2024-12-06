FROM golang:1.22.6 as base

WORKDIR /go/src/work.ctyun.cn/git/cwai/cwai-event-service
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o  cwai-event-service  cmd/server/server.go

FROM alpine:3.19

COPY --from=base /go/src/work.ctyun.cn/git/cwai/cwai-event-service/cwai-event-service /bin/

WORKDIR /bin
CMD ["/bin/cwai-event-service","-c","/etc/cwai-event-service/config.yaml"]