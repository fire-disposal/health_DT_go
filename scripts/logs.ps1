# 查看 server.exe 的最新日志（假设有日志文件 server.log）
$OutputEncoding = [Console]::OutputEncoding = [Text.Encoding]::UTF8
$logPath = "$PSScriptRoot\..\server.log"
if (Test-Path $logPath) {
    Get-Content $logPath -Tail 50
} else {
    Write-Host "未找到日志文件 $logPath" -Encoding UTF8
}