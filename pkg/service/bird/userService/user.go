package userService

import (
	"go-trailer-api/pkg/model/bird/userModel"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"` //用户名
	Password string `form:"password" json:"password" binding:"required"` //密码
	Nickname string `form:"nickname" json:"nickname" binding:"required"` //昵称
	Gender   int    `form:"gender" json:"gender" binding:"gender"`       //性别
	Birthday string `form:"birthday" json:"birthday" binding:""`
}

func MapUser(user User) map[string]interface{} {
	return map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"nickname": user.Nickname,
		"gender":   user.Gender,
		"birthday": user.Birthday,
	}
}

//新增
func (u User) Add() error {
	userInfo := MapUser(u)

	if err := userModel.AddUser(userInfo); err != nil {
		return err
	}

	return nil
}

//获取列表
func Listing(params userModel.UserListParams) (error, list interface{}, total int64) {

	return userModel.UserList(params)
}
