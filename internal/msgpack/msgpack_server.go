// Package msgpack 提供独立的 TCP 监听与协议解析能力
package msgpack

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
)

// MAGIC 头与端口定义
const (
	TCPPort = 5858
	MAGIC   = "\xab\xcd"
)

// CRC 校验
func crc8(data []byte) byte {
	var crc byte = 0
	for _, b := range data {
		crc ^= b
		for i := 0; i < 8; i++ {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ 0x31
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

// PayloadHandler 业务处理回调
type PayloadHandler func(payload map[string]interface{})

// MsgpackServer 结构体
type MsgpackServer struct {
	handler PayloadHandler
}

// NewMsgpackServer 构造
func NewMsgpackServer(handler PayloadHandler) *MsgpackServer {
	return &MsgpackServer{handler: handler}
}

// Start 启动监听
func (s *MsgpackServer) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", TCPPort))
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		go s.handleConn(conn)
	}
}

// handleConn 处理单连接
func (s *MsgpackServer) handleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 0)
	tmp := make([]byte, 1024)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				zap.L().Error("TCP连接读取异常", zap.Error(err))
			}
			break
		}
		buffer = append(buffer, tmp[:n]...)
		for len(buffer) >= 4 {
			if !bytes.HasPrefix(buffer, []byte(MAGIC)) {
				idx := bytes.Index(buffer, []byte(MAGIC))
				if idx == -1 {
					buffer = buffer[:0]
					break
				}
				buffer = buffer[idx:]
				if len(buffer) < 4 {
					break
				}
			}
			length := buffer[2]
			crcVal := buffer[3]
			if length < 1 {
				buffer = buffer[4:]
				continue
			}
			if len(buffer) < int(4+length) {
				break
			}
			data := buffer[4 : 4+length]
			if crc8(data) != crcVal {
				buffer = buffer[4+length:]
				continue
			}
			// 业务解包交由 handler
			if s.handler != nil {
				var unpacked map[string]interface{}
				if err := msgpack.Unmarshal(data, &unpacked); err == nil {
					s.handler(unpacked)
				}
			}
			buffer = buffer[4+length:]
		}
	}
}
