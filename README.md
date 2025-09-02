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
1. 推荐使用 `config/config.yaml` 进行集中配置（支持环境变量自动兼容）。
2. 安装依赖：`go get github.com/spf13/viper`
3. 运行 `go run cmd/server/main.go`
4. 参考 `api/http/` 目录进行接口开发

### 配置方式说明

- 优先读取 `config/config.yaml`，如未找到则自动回退环境变量。
- 配置项与环境变量一一对应，详见 [`config/env.example`](config/env.example:1)。

#### config.yaml 示例

```yaml
server:
  port: 8002
  msglistener_port: 5858
postgres:
  host: localhost
  port: 5432
  user: postgres
  password: 12345678
  dbname: health_dt
  sslmode: disable
redis:
  addr: localhost:6379
  password: ""
  db: 0
mqtt:
  broker: tcp://localhost:1883
  client_id: health_dt_client
  username: ""
  password: ""
websocket:
  host: 0.0.0.0
  port: 8765
  path: /ws/health
jwt_secret: your-secret-key-please-change-in-production
wechat:
  appid: your-wechat-appid
  secret: your-wechat-secret
```

### 常见问题排查

- 配置加载失败：请检查 `config/config.yaml` 路径及格式，或环境变量是否正确设置。
- 数据库认证失败：请确认密码、端口、用户名与数据库实际一致，查看日志详细输出定位问题。
- 依赖缺失：请确保已执行 `go get github.com/spf13/viper`。

更多细节见 [`docs/plan_ai.md`](docs/plan_ai.md:1)