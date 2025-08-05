#!/bin/bash

# 编译
echo "Building application..."
go build -o recruit-center

# 构建Docker镜像
echo "Building Docker image..."
docker-compose build

echo "Build completed successfully!"
