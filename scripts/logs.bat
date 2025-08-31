@echo off
chcp 65001 >nul
set "logPath=..\server.log"
if exist "%logPath%" (
    powershell -Command "Get-Content '%logPath%' -Tail 50"
) else (
    echo 未找到日志文件 %logPath%
)