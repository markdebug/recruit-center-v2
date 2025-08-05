#!/bin/bash

# 停止现有服务
echo "Stopping existing services..."
docker-compose down

# 启动服务
echo "Starting services..."
docker-compose up -d

echo "Deployment completed successfully!"
