# 运行 server.exe
$OutputEncoding = [Console]::OutputEncoding = [Text.Encoding]::UTF8
Write-Host "启动 server.exe ..." -Encoding UTF8
Start-Process -NoNewWindow -FilePath "$PSScriptRoot\..\server.exe"