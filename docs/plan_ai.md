# 健康数据平台（health_DT_go）AI友好型开发方案

## 1. 技术栈推荐

- Web/API：Gin v1.10.1
- 消息队列：paho.mqtt.golang v1.5.0
- 配置管理：viper v1.20.1
- 数据库驱动：lib/pq v1.10.9
- JWT认证：golang-jwt/jwt/v4 v4.5.2
- 加密工具：x/crypto v0.41.0
- 日志：go.uber.org/zap v1.27.0
- 其他：goroutines/channels

---

## 2. 标准化目录结构

```plaintext
```plaintext
health_DT_go/
├─ cmd/
│  └─ server/
│     └─ main.go
├─ config/
│  ├─ config.yaml
│  └─ config.go
├─ internal/
│  ├─ app/
│  │  ├─ eventbus/
│  │  ├─ handlers/
│  │  │   ├─ health/
│  │  │   │   ├─ base.go
│  │  │   │   ├─ heart_rate_handler.go
│  │  │   │   ├─ blood_pressure_handler.go
│  │  │   │   ├─ spo2_handler.go
│  │  │   │   └─ temperature_handler.go
│  │  └─ pipeline.go
│  ├─ models/
│  ├─ repository/
│  │  └─ postgres/
│  │     ├─ auth_repo.go
│  │     └─ user_repo.go
│  ├─ service/
│  ├─ ws/
│  └─ mqtt/
│     └─ mqtt_client.go
├─ api/
│  ├─ http/
│  │  ├─ auth_routes.go
│  │  ├─ middleware.go
│  │  └─ user_routes.go
│  └─ grpc/
├─ pkg/
│  ├─ logger/
│  ├─ utils/
│  │  └─ password.go
│  └─ error/
├─ scripts/
├─ migrations/
├─ go.mod
└─ go.sum
```

**鉴权流程：**
- 采用JWT进行Token认证
- 通过中间件校验接口权限
- 区分admin与app角色，便于权限管理

---

## 5. 系统特点与扩展性

- 事件驱动：设备数据、告警、推送均通过事件总线流转
- 支持模拟数据：可选择不入库，仅推送或缓存
- 异步批量处理：goroutines与队列提升入库性能
- 前后端解耦：WebSocket与HTTP接口独立部署
- Redis高并发：缓存最新状态与模拟数据，减轻数据库压力
- 高可扩展性：新增设备、事件、告警规则仅需添加Handler，无需改动核心流程

---

## 7. 建议优化（结合 Redis）

### repository 分层

当前只有 `repository/postgres/`，建议增加一层：

```plaintext
internal/repository/
├─ postgres/
├─ redis/
└─ interfaces.go   // 定义通用接口
```

这样业务逻辑可灵活选择 PostgreSQL 持久化或 Redis 临时存储。

### Redis 的定位

建议拆分两类用途：

- **状态缓存**：保存设备/用户的“最新一份”数据，随时覆盖（SET device:{id}:state ...）。
- **短期流数据**：用 Redis Stream 或 List+LTRIM 保存最近 N 条消息（如 1 万条），自动丢弃旧数据。

### 数据流向优化

当前 Mermaid 架构图 D[事件处理器 handlers] → E[数据入库 repository] 直接入 PostgreSQL。建议增加 Redis 分支：

```mermaid
graph TD
    A[设备/模拟数据] --> B[MQTT/HTTP 接入]
    B --> C[事件总线 eventbus]
    C --> D[事件处理器 handlers]
    D --> E[数据入库 repository(PostgreSQL)]
    D --> F[Redis缓存/Stream]
    D --> G[告警生成]
    D --> H[WebSocket 推送]
    E --> I[PostgreSQL]
    F --> H
    G --> J[告警通知]
    H --> K[前端]
```

👉 Redis 主要做两件事：

- 承接无意义或短期展示数据，不入库。
- 给 WebSocket 提供低延迟数据源。

### Go 代码层建议

定义 DataSink 接口：

```go
type DataSink interface {
    Save(ctx context.Context, data *models.HealthData) error
}
```

- `postgresSink` 实现长久存储
- `redisSink` 实现临时存储

handler 可根据数据来源（模拟/现场 vs. 真正业务）选择 sink。
## 6. Mermaid 架构图（简化版）

```mermaid
graph TD
    A[设备/模拟数据] --> B[MQTT/HTTP 接入]
    B --> C[事件总线 eventbus]
    C --> D[事件处理器 handlers]
    D --> E[数据入库 repository]
    D --> F[告警生成]
    D --> G[WebSocket 推送]
    E --> H[PostgreSQL]
    G --> I[前端]
    F --> J[告警通知]
---
