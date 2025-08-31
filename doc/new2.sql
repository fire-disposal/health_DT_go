-- ================================================
-- 用户与登录相关（未改动）
-- ================================================

CREATE TABLE admin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(128),
    phone VARCHAR(32),
    password_hash VARCHAR(256) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE app_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(128),
    phone VARCHAR(32),
    password_hash VARCHAR(256) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP,
    wechat_openid VARCHAR(128),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ================================================
-- 核心业务表
-- ================================================

-- 健康档案表（HealthProfile） 保留原结构
CREATE TABLE health_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES app_users(id) ON DELETE SET NULL,
    name VARCHAR(128) NOT NULL,
    gender VARCHAR(16),
    birth_date DATE,
    metadata JSONB,  -- 保留灵活字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_health_profiles_user ON health_profiles(user_id);

-- ----------------------------
-- 设备表（devices） 极简化
-- ----------------------------
CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    serial_number VARCHAR(64) UNIQUE NOT NULL,  -- 新增唯一标识
    name VARCHAR(128),                          -- 保留友好名
    device_type VARCHAR(64),                    -- 类型区分
    is_active BOOLEAN DEFAULT TRUE,             -- 激活状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_devices_type ON devices(device_type);

-- ----------------------------
-- 设备与健康档案绑定表（新增）
-- ----------------------------
CREATE TABLE device_assignments (
    id SERIAL PRIMARY KEY,
    device_id INT REFERENCES devices(id) ON DELETE CASCADE,
    health_profile_id INT REFERENCES health_profiles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    unassigned_at TIMESTAMP
);
CREATE INDEX idx_device_assignments_profile ON device_assignments(health_profile_id);
CREATE INDEX idx_device_assignments_device ON device_assignments(device_id);

-- ----------------------------
-- 健康数据记录表（health_data_records） 保留
-- ----------------------------
CREATE TABLE health_data_records (
    id SERIAL PRIMARY KEY,
    health_profile_id INT REFERENCES health_profiles(id) ON DELETE CASCADE,
    device_id INT REFERENCES devices(id) ON DELETE SET NULL,
    schema_type VARCHAR(64),
    recorded_at TIMESTAMP NOT NULL,
    payload JSONB,  -- 数据本身用 JSONB 保留灵活性
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_hdata_profile_time ON health_data_records(health_profile_id, recorded_at);
CREATE INDEX idx_hdata_device_time ON health_data_records(device_id, recorded_at);

-- ----------------------------
-- 事件表（events） 保留
-- ----------------------------
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(64) NOT NULL,
    health_profile_id INT REFERENCES health_profiles(id) ON DELETE CASCADE,
    device_id INT REFERENCES devices(id) ON DELETE SET NULL,
    source_record_id INT REFERENCES health_data_records(id) ON DELETE SET NULL,
    timestamp TIMESTAMP NOT NULL,
    data JSONB,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_events_type_time ON events(event_type, timestamp);
CREATE INDEX idx_events_profile_time ON events(health_profile_id, timestamp);
CREATE INDEX idx_events_device_time ON events(device_id, timestamp);

-- ----------------------------
-- 告警表（alerts） 保留
-- ----------------------------
CREATE TABLE alerts (
    id SERIAL PRIMARY KEY,
    health_profile_id INT REFERENCES health_profiles(id) ON DELETE SET NULL,
    device_id INT REFERENCES devices(id) ON DELETE SET NULL,
    source_event_id INT REFERENCES events(id) ON DELETE SET NULL,
    rule_name VARCHAR(128),
    level VARCHAR(32),
    message TEXT,
    event_type VARCHAR(64),
    description TEXT,
    extra JSONB,
    status VARCHAR(32),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP
);
CREATE INDEX idx_alerts_device_status ON alerts(device_id, status);
CREATE INDEX idx_alerts_profile_rule_status ON alerts(health_profile_id, rule_name, status);
