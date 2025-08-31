# 清理编译生成的文件
$OutputEncoding = [Console]::OutputEncoding = [Text.Encoding]::UTF8
Write-Host "正在清理 server.exe ..." -Encoding UTF8
Remove-Item "$PSScriptRoot\..\server.exe" -ErrorAction SilentlyContinue
Write-Host "清理完成。" -Encoding UTF8