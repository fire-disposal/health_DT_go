// 用户相关路由
package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fire-disposal/health_DT_go/internal/models"
	"github.com/fire-disposal/health_DT_go/internal/service"
)

/*
// @Summary 获取单个用户信息
// @Description 根据ID查询用户信息
// @Tags user
// @Produce json
// @Param id query int true "用户ID"
// @Success 200 {object} models.AppUser
// @Failure 404 {string} string "用户不存在"
// @Router /user/info [get]
*/
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	id := int64(idInt)
	user, err := service.GetUserByID(id)
	if err != nil {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

/*
// @Summary 获取用户列表
// @Description 查询所有用户信息
// @Tags user
// @Produce json
// @Success 200 {array} models.AppUser
// @Router /user/list [get]
*/
func UserListHandler(w http.ResponseWriter, r *http.Request) {
	users, err := service.ListUsers()
	if err != nil {
		http.Error(w, "获取用户列表失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

/*
// @Summary 用户注册
// @Description 新用户注册
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.AppUser true "用户信息"
// @Success 200 {object} models.AppUser
// @Failure 500 {string} string "注册失败"
// @Router /user/register [post]
*/
func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AppUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	user, err := service.CreateUser(&req)
	if err != nil {
		http.Error(w, "注册失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

/*
// @Summary 用户信息更新
// @Description 更新用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.AppUser true "用户信息"
// @Success 200 {object} models.AppUser
// @Failure 500 {string} string "更新失败"
// @Router /user/update [put]
*/
func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AppUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == 0 {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	user, err := service.UpdateUser(&req)
	if err != nil {
		http.Error(w, "更新失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
