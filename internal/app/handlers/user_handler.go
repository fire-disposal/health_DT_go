// 用户事件处理器，仅处理用户相关事件，结构便于扩展
package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/service"
)

// UserHandler 负责处理用户相关事件
type UserHandler struct {
	UserService service.UserService
}

// NewUserHandler 构造函数，便于依赖注入和扩展
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// Register 处理用户注册事件
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败", http.StatusBadRequest)
		return
	}
	user := &models.AppUser{Username: req.Username}
	err := h.UserService.Register(r.Context(), user, req.Password)
	if err != nil {
		http.Error(w, "注册失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetInfo 处理用户信息查询事件
func (h *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	// 通过 query 参数 user_id 查询
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "缺少 user_id 参数", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "user_id 参数格式错误", http.StatusBadRequest)
		return
	}
	user, err := h.UserService.GetUserInfo(r.Context(), userID)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
