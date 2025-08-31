# 健康数据平台（health_DT_go）

一个基于 Go 的健康数据采集与处理平台，支持设备数据、告警、推送等事件驱动流程，易于扩展与二次开发。

## 技术栈
- Gin（Web/API）
- paho.mqtt.golang（消息队列）
- viper（配置管理）
- lib/pq（PostgreSQL驱动）
- golang-jwt/jwt（认证）
- goroutines/channels（并发）

## 目录结构
```
health_DT_go/
├─ cmd/server/main.go
├─ config/config.yaml
├─ internal/app/handlers/health/
├─ internal/models/
├─ internal/repository/postgres/
├─ internal/service/
├─ api/http/
├─ scripts/
├─ docs/
├─ go.mod
└─ go.sum
```

## 系统特点
- 事件总线驱动，支持设备数据、告警、推送
- JWT鉴权，支持多角色权限管理
- 支持模拟数据与异步批量处理
- 前后端解耦，支持高并发缓存
- 易扩展：新增设备/事件/告警仅需添加 Handler

## 快速开始
1. 配置数据库与环境变量
2. 运行 `go run cmd/server/main.go`
3. 参考 `api/http/` 目录进行接口开发

更多细节见 [`docs/plan_ai.md`](docs/plan_ai.md:1)