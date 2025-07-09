# 多阶段构建
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# 最终镜像
FROM alpine:latest

# 安装ca-certificates用于HTTPS请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从builder阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY .env .env

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]