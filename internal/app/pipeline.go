// internal/app/pipeline.go
package app

import (
	"github.com/fire-disposal/health_DT_go/internal/app/eventbus"
)

// HealthEvent 统一健康数据事件结构体
type HealthEvent struct {
	DeviceID  string      // 设备ID或模拟标识
	EventType string      // 事件类型：heart_rate, blood_pressure, spo2, temperature
	Payload   interface{} // 具体数据载体
	Source    string      // 来源标识（设备/模拟）
}

// HealthDataProcessor 健康数据处理器接口，便于扩展
type HealthDataProcessor interface {
	Handle(event HealthEvent)
}

// Pipeline 健康数据处理主流程
type Pipeline struct {
	processors []HealthDataProcessor
	eventBus   *eventbus.EventBus
}

// NewPipeline 创建主流程实例
func NewPipeline(bus *eventbus.EventBus) *Pipeline {
	return &Pipeline{
		processors: make([]HealthDataProcessor, 0),
		eventBus:   bus,
	}
}

// RegisterProcessor 注册健康数据处理器，支持扩展
func (p *Pipeline) RegisterProcessor(processor HealthDataProcessor) {
	p.processors = append(p.processors, processor)
}

// ReceiveEvent 统一接收事件并分发
func (p *Pipeline) ReceiveEvent(event HealthEvent) {
	// 分发至 eventbus
	if p.eventBus != nil {
		p.eventBus.Publish(event.EventType, event)
	}
	// 分发至各健康数据处理器
	for _, processor := range p.processors {
		processor.Handle(event)
	}
}
