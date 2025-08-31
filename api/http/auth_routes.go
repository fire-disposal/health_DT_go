// 鉴权相关路由
package http

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/fire-disposal/health_DT_go/internal/auth"
	"github.com/fire-disposal/health_DT_go/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		LoginType string `json:"login_type"` // "admin" 或 "app"
		Username  string `json:"username"`
		Password  string `json:"password"`
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=health_dt sslmode=disable")
	if err != nil {
		http.Error(w, "数据库连接失败", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	repo := postgres.NewAuthRepository(db)

	var userID int64
	var role string

	if req.LoginType == "admin" {
		admin, err := repo.GetAdminUserByUsername(req.Username)
		if err != nil || admin == nil {
			http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
			http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
			return
		}
		userID = admin.ID
		role = admin.Role
	} else if req.LoginType == "app" {
		user, err := repo.GetAppUserByUsername(req.Username)
		if err != nil || user == nil {
			http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
			return
		}
		userID = user.ID
		role = "app"
	} else {
		http.Error(w, "未知登录类型", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(userID, role)
	if err != nil {
		http.Error(w, "生成Token失败", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 管理员登录（带 db 参数）
func AdminLoginHandlerWithDB(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	repo := postgres.NewAuthRepository(db)
	admin, err := repo.GetAdminUserByUsername(req.Username)
	if err != nil || admin == nil {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateToken(admin.ID, admin.Role)
	if err != nil {
		http.Error(w, "生成Token失败", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 普通用户登录（带 db 参数）
func AppLoginHandlerWithDB(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	repo := postgres.NewAuthRepository(db)
	user, err := repo.GetAppUserByUsername(req.Username)
	if err != nil || user == nil {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateToken(user.ID, "app")
	if err != nil {
		http.Error(w, "生成Token失败", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
