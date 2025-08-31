// Package health 提供健康数据处理器基类接口及通用方法定义。
package health

import (
	"context"
)

// HealthEvent 表示健康相关的事件数据结构。
// 可根据实际需求扩展字段。
type HealthEvent struct {
	Type string      // 事件类型
	Data interface{} // 事件数据
}

// HealthHandler 健康数据处理器接口，定义通用方法。
type HealthHandler interface {
	// HandleEvent 处理健康事件，返回处理结果或错误。
	HandleEvent(ctx context.Context, event HealthEvent) error

	// ValidateData 校验健康数据的有效性，返回校验结果或错误。
	ValidateData(data interface{}) error
}

// BaseHealthHandler 提供健康处理器通用方法的基础实现。
// 可嵌入具体处理器以复用通用逻辑。
type BaseHealthHandler struct{}

// ValidateData 默认实现，需具体处理器重写。
func (b *BaseHealthHandler) ValidateData(data interface{}) error {
	// 默认不做校验，直接通过。
	return nil
}

// HandleEvent 默认实现，需具体处理器重写。
func (b *BaseHealthHandler) HandleEvent(ctx context.Context, event HealthEvent) error {
	// 默认不做处理，直接通过。
	return nil
}
