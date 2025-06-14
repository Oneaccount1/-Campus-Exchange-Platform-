@echo off
REM 检查是否存在测试数据库
echo 检查测试数据库是否存在...

REM 创建测试数据库（如果不存在）
mysql -u root -ppassword -e "CREATE DATABASE IF NOT EXISTS campus_exchange_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo 测试数据库已准备就绪

REM 设置环境变量以启用真实数据库测试
set RUN_REAL_DB_TESTS=true

REM 运行指定的测试文件
echo 运行真实数据库测试...
go test -v ./internal/modules/user/repositories -run TestRealDatabaseSuite

REM 运行完成后清理测试数据（可选）
REM echo 清理测试数据库...
REM mysql -u root -ppassword -e "DROP DATABASE campus_exchange_test;"

echo 测试完成
pause 