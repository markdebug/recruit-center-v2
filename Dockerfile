FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o recruit-center .

# 使用轻量级基础镜像
FROM alpine:latest

WORKDIR /app

# 复制配置文件
COPY config.yaml .
# 从builder阶段复制编译好的应用
COPY --from=builder /app/recruit-center .

EXPOSE 8080

CMD ["./recruit-center"]
