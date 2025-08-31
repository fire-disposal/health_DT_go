# 编译 Go 项目为 server.exe
Write-Host "正在编译 server.exe ..." -Encoding UTF8
go build -o server.exe ./cmd/server/main.go
if ($LASTEXITCODE -eq 0) {
    Write-Host "编译成功！" -Encoding UTF8
} else {
    Write-Host "编译失败！" -Encoding UTF8
}