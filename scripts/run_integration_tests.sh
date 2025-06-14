#!/bin/bash

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
  echo "错误：Docker未运行或无法连接到Docker守护进程"
  exit 1
fi

# 设置环境变量以启用集成测试
export RUN_INTEGRATION_TESTS=true

# 运行集成测试
echo "运行集成测试..."
go test -v -tags=integration ./internal/modules/user/repositories -run TestIntegrationSuite

echo "测试完成" 