# 健康数据平台（health_DT_go）开发方案

## 1. 技术栈

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
├─ cmd/                  # 项目入口，服务启动
│  └─ server/
│     └─ main.go         # 主程序入口
├─ config/               # 配置管理
│  ├─ config.go          # 配置加载逻辑
│  └─ env.example        # 环境变量示例
├─ internal/
│  ├─ app/               # 核心应用逻辑
│  │  ├─ eventbus/       # 事件驱动总线
│  │  │   └─ eventbus.go # 事件分发实现
│  │  ├─ handlers/       # 业务处理器
│  │  │   ├─ auth_handler.go           # 认证处理
│  │  │   ├─ mqtt_handler.go           # MQTT数据处理
│  │  │   ├─ msgpack_event_handler.go  # MsgPack事件处理
│  │  │   ├─ user_handler.go           # 用户业务处理
│  │  │   └─ health/                   # 健康数据处理
│  │  │       ├─ base.go
│  │  │       ├─ heart_rate_handler.go
│  │  │       ├─ blood_pressure_handler.go
│  │  │       ├─ spo2_handler.go
│  │  │       └─ temperature_handler.go
│  │  └─ pipeline.go    # 数据流管道
│  ├─ models/           # 数据结构定义
│  │  ├─ admin_user.go
│  │  ├─ alerts.go
│  │  ├─ app_user.go
│  │  ├─ auth.go
│  │  ├─ device_assignments.go
│  │  ├─ devices.go
│  │  ├─ events.go
│  │  ├─ health_data_records.go
│  │  └─ health_profiles.go
│  ├─ repository/       # 数据持久化
│  │  └─ postgres/
│  │     ├─ alerts_repo.go        # 告警数据存储
│  │     ├─ auth_repo.go          # 认证数据存储
│  │     ├─ events_repo.go        # 事件数据存储
│  │     ├─ health_data_repo.go   # 健康数据存储
│  │     └─ user_repo.go          # 用户数据存储
│  ├─ service/          # 业务服务层
│  │  ├─ auth_service.go         # 认证服务
│  │  └─ user_service.go         # 用户服务
│  ├─ mqtt/             # MQTT客户端
│  │  └─ mqtt_client.go
│  ├─ msgpack/          # MsgPack服务端
│  │  └─ msgpack_server.go
├─ api/
│  ├─ http/             # RESTful 路由
│  │  ├─ alerts_routes.go    # 告警接口
│  │  ├─ auth_routes.go      # 认证接口
│  │  ├─ events_routes.go    # 事件接口
│  │  ├─ health_routes.go    # 健康数据接口
│  │  ├─ middleware.go       # 路由中间件
│  │  └─ user_routes.go      # 用户接口
├─ docs/                # 项目文档
│  ├─ docs.go
│  ├─ new3.sql
│  ├─ plan_ai.md
│  ├─ swagger.json
│  └─ swagger.yaml
├─ pkg/                 # 工具包
│  ├─ logger/           # 日志工具
│  ├─ utils/            # 通用工具
│  │  └─ password.go    # 密码工具
│  └─ error/            # 错误处理
├─ scripts/             # 辅助脚本
├─ migrations/          # 数据库迁移
├─ go.mod               # 依赖声明
└─ go.sum               # 依赖校验
```

**鉴权流程：**
- 采用JWT进行Token认证
- 通过中间件校验接口权限
- 区分admin与app角色，便于权限管理

---

## 项目特性与优化建议（实现状态一览）

- [x] 事件驱动：设备数据、告警、推送均通过事件总线流转
- [x] 支持模拟数据：可选择不入库，仅推送或缓存
- [x] 异步批量处理：goroutines与队列提升入库性能
- [x] 前后端解耦：WebSocket与HTTP接口独立部署
- [ ] Redis高并发：缓存最新状态与模拟数据，减轻数据库压力
- [x] 高可扩展性：新增设备、事件、告警规则仅需添加Handler，无需改动核心流程
- [ ] repository 分层：建议增加 redis 层及通用接口，灵活选择存储方案
- [ ] Redis状态缓存/流数据：建议拆分状态缓存与短期流数据，提升实时性
- [ ] 数据流向优化：建议在事件处理器后增加 Redis 分支，支持缓存与流式推送
- [ ] DataSink接口：建议定义通用数据存储接口，便于扩展多种存储后端

> 已实现项已打勾，未实现项留空，便于团队与 AI 自动化跟踪进度。
## 6. Mermaid 架构图（简化版，同步实际路由与数据流）

```mermaid
graph TD
    A[设备/模拟数据] --> B[MQTT/HTTP 接入]
    B --> C[事件总线 eventbus]
    C --> D[事件处理器 handlers]
    D --> E[数据入库 repository(postgres)]
    D --> F[告警生成]
    D --> G[WebSocket 推送]
    D --> H[HTTP/REST API]
    D --> I[Events API]
    D --> J[Health API]
    E --> K[PostgreSQL]
    F --> L[告警通知]
    G --> M[前端]
    H --> M
    I --> M
    J --> M
```

