# 使用官方的 Golang 镜像，指定版本为 1.23.2。这里已经默认是 linux/amd64 的镜像。
FROM golang:1.22.2

# 设置工作目录
WORKDIR /app

# 将当前目录下的所有文件复制到工作目录中
COPY . .

# 设置环境变量并编译
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go build -o main .

# 暴露端口号 8081
EXPOSE 8081

# 设置容器启动时运行的命令
CMD ["./main"]