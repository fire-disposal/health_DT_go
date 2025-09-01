// generator.go
package simdata

import (
	"math/rand"
	"time"
)

// 类型定义（可后续抽象到 models 层）
type HeartRateEventData struct {
	UserID    string
	HeartRate int
	Timestamp int64
}

type BloodPressureEventData struct {
	UserID    string
	Systolic  int
	Diastolic int
	Timestamp int64
}

type SpO2EventData struct {
	UserID    string
	SpO2      int
	Timestamp int64
}

type TemperatureEventData struct {
	UserID      string
	Temperature float64
	Timestamp   int64
}

// 心率模拟
func GenerateHeartRate(userID string) HeartRateEventData {
	return HeartRateEventData{
		UserID:    userID,
		HeartRate: rand.Intn(40) + 60,
		Timestamp: time.Now().Unix(),
	}
}

// 血压模拟
func GenerateBloodPressure(userID string) BloodPressureEventData {
	return BloodPressureEventData{
		UserID:    userID,
		Systolic:  rand.Intn(60) + 90,
		Diastolic: rand.Intn(40) + 60,
		Timestamp: time.Now().Unix(),
	}
}

// 血氧模拟
func GenerateSpO2(userID string) SpO2EventData {
	return SpO2EventData{
		UserID:    userID,
		SpO2:      rand.Intn(6) + 95,
		Timestamp: time.Now().Unix(),
	}
}

// 体温模拟
func GenerateTemperature(userID string) TemperatureEventData {
	return TemperatureEventData{
		UserID:      userID,
		Temperature: rand.Float64()*2 + 36.0,
		Timestamp:   time.Now().Unix(),
	}
}
