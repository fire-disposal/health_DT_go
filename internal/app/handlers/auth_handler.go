// Package handlers 实现鉴权事件处理器
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fire-disposal/health_DT_go/internal/service"
)

// AuthHandlerInterface 便于扩展的鉴权事件处理器接口
type AuthHandlerInterface interface {
	Login(w http.ResponseWriter, r *http.Request)
	ValidateToken(w http.ResponseWriter, r *http.Request)
}

// AuthHandler 实现 AuthHandlerInterface
type AuthHandler struct {
	AuthService *service.AuthService
}

// Login 处理登录事件，生成 Token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	token, err := h.AuthService.GenerateToken(req.UserID, 24*time.Hour)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ValidateToken 处理 Token 校验事件
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	auth, err := h.AuthService.ValidateToken(req.Token)
	if err != nil {
		http.Error(w, "token invalid or expired", http.StatusUnauthorized)
		return
	}
	resp := map[string]interface{}{
		"user_id":   auth.UserID,
		"expiresAt": auth.ExpiresAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
