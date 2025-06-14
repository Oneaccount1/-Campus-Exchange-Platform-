#!/bin/bash

# 检查是否存在测试数据库
echo "检查测试数据库是否存在..."
MYSQL_CMD="mysql -u root -ppassword -e"

# 创建测试数据库（如果不存在）
$MYSQL_CMD "CREATE DATABASE IF NOT EXISTS campus_exchange_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo "测试数据库已准备就绪"

# 设置环境变量以启用真实数据库测试
export RUN_REAL_DB_TESTS=true

# 运行指定的测试文件
echo "运行真实数据库测试..."
go test -v ./internal/modules/user/repositories -run TestRealDatabaseSuite

# 运行完成后清理测试数据（可选）
# echo "清理测试数据库..."
# $MYSQL_CMD "DROP DATABASE campus_exchange_test;"

echo "测试完成" 