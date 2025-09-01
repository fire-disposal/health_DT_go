
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
│  │  └─ pipeline.go    # 健康数据主流程：统一事件分发，支持多处理器扩展
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
│  │  ├─ postgres/
│  │  │   ├─ alerts_repo.go            # 告警数据存储
│  │  │   ├─ auth_repo.go              # 认证数据存储
│  │  │   ├─ devices_repo.go           # 设备数据存储
│  │  │   ├─ events_repo.go            # 事件数据存储
│  │  │   ├─ health_data_repo.go       # 健康数据存储
│  │  │   ├─ health_profiles_repo.go   # 健康档案存储
│  │  │   └─ user_repo.go              # 用户数据存储
│  │  ├─ redis/
│  │  │   ├─ redis_client.go           # Redis客户端
│  │  │   └─ simdata_repo.go           # 模拟数据存储
│  ├─ service/          # 业务服务层
│  │  ├─ auth_service.go               # 认证服务
│  │  ├─ devices_service.go            # 设备服务
│  │  ├─ health_profiles_service.go    # 健康档案服务
│  │  └─ user_service.go               # 用户服务
│  ├─ mqtt/             # MQTT客户端
│  │  └─ mqtt_client.go
│  ├─ msgpack/          # MsgPack服务端
│  │  └─ msgpack_server.go
│  ├─ simdata/          # 数据模拟
│  │  └─ generator.go
├─ api/
│  ├─ http/             # RESTful 路由
│  │  ├─ alerts_routes.go            # 告警接口
│  │  ├─ auth_routes.go              # 认证接口
│  │  ├─ devices_routes.go           # 设备接口
│  │  ├─ events_routes.go            # 事件接口
│  │  ├─ health_profiles_routes.go   # 健康档案接口
│  │  ├─ health_routes.go            # 健康数据接口
│  │  ├─ middleware.go               # 路由中间件
│  │  └─ user_routes.go              # 用户接口
├─ docs/                # 项目文档
│  ├─ docs.go
│  ├─ new3.sql
│  ├─ plan_ai.md
│  ├─ swagger.json
│  └─ swagger.yaml
├─ scripts/             # 辅助脚本
├─ migrations/          # 数据库迁移
├─ go.mod               # 依赖声明
├─ go.sum               # 依赖校验
```