FROM golang:1.22.0
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 8081
CMD ["sh", "-ce", "./main"]
