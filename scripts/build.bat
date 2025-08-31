@echo off
chcp 65001 >nul
echo 正在编译 server.exe ...
go build -o server.exe ./cmd/server/main.go
if %errorlevel%==0 (
    echo 编译成功！
) else (
    echo 编译失败！
)