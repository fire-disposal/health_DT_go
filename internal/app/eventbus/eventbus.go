// Package eventbus 实现一个简单且可扩展的事件总线。
package eventbus

import (
	"sync"
)

// SimDataEvent 模拟数据事件结构
type SimDataEvent struct {
	EventType string      // heart_rate, blood_pressure, etc.
	Payload   interface{} // 具体数据结构
	Source    string      // 来源标识（simulator）
	Timestamp int64       // 事件时间戳
}

// EventHandler 事件处理函数类型
type EventHandler func(data any)

// EventBus 事件总线结构体
type EventBus struct {
	mu        sync.RWMutex
	listeners map[string][]EventHandler
}

// NewEventBus 创建一个新的事件总线实例
func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅某类型事件
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	// 允许重复订阅，直接追加
	eb.listeners[eventType] = append(eb.listeners[eventType], handler)
}

// Unsubscribe 取消订阅某类型事件
func (eb *EventBus) Unsubscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if handlers, ok := eb.listeners[eventType]; ok {
		for i, h := range handlers {
			// 只移除第一个指针相同的 handler
			if &h == &handler {
				eb.listeners[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
		if len(eb.listeners[eventType]) == 0 {
			delete(eb.listeners, eventType)
		}
	}
}

// Publish 发布事件
func (eb *EventBus) Publish(eventType string, data any) {
	eb.mu.RLock()
	handlers := eb.listeners[eventType]
	eb.mu.RUnlock()
	for _, handler := range handlers {
		go handler(data)
	}
}
