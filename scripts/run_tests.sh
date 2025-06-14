#!/bin/bash

# 运行所有测试并生成覆盖率报告
go test -v ./... -coverprofile=coverage.out

# 显示覆盖率报告
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report generated at coverage.html"

# 运行特定测试
# go test -v ./internal/modules/user/repositories -run TestUserRepositorySuite 