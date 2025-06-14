@echo off

REM 检查Docker是否运行
docker info > nul 2>&1
if %ERRORLEVEL% NEQ 0 (
  echo 错误：Docker未运行或无法连接到Docker守护进程
  exit /b 1
)

REM 设置环境变量以启用集成测试
set RUN_INTEGRATION_TESTS=true

REM 运行集成测试
echo 运行集成测试...
go test -v -tags=integration ./internal/modules/user/repositories -run TestIntegrationSuite

echo 测试完成
pause 