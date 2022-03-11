package userModel

import (
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/model"
	"go-trailer-api/pkg/util"
	"strconv"
)

type Users struct {
	Id         int    `json:"id" gorm:"column:id"`
	Username   string `json:"username" gorm:"column:username"`       //用户名
	Password   string `json:"password" gorm:"column:password"`       //密码
	Salt       string `json:"salt" gorm:"column:salt"`               //密码盐
	Nickname   string `json:"nickname" gorm:"column:nickname"`       //昵称
	Gender     int    `json:"gender" gorm:"column:gender"`           //性别 - 0：男；1：女
	Birthday   string `json:"birthday" gorm:"column:birthday"`       //生日
	CreateTime string `json:"create_time" gorm:"column:create_time"` //创建时间
}

type UserListParams struct {
	Page     int    `form:"page"`                       //页码
	PageSize int    `form:"page_size" json:"page_size"` //每页数量
	Username string `form:"username"`                   //用户名
	Nickname string `form:"nickname"`                   //昵称
}

var tableName = "go_users"

func (Users) TableName() string {
	return tableName
}

func AddUser(data map[string]interface{}) error {
	salt := strconv.FormatInt(util.RandIntNumber(100000, 999999), 10)
	pwd := util.Md5V(data["password"].(string) + salt)

	user := Users{
		Username:   data["username"].(string),
		Password:   pwd,
		Salt:       salt,
		Nickname:   data["nickname"].(string),
		Gender:     data["gender"].(int),
		Birthday:   data["birthday"].(string),
		CreateTime: util.GetCurrentTime(),
	}

	if err := model.BirdDb.Create(&user).Error; err != nil {
		logging.Error(err)

		return err
	}

	return nil
}

func UserList(params UserListParams) (error, list interface{}, total int64) {
	var lists []Users

	db := model.BirdDb.Table(tableName)
	if params.Username != "" {
		db = db.Where("username like ? ", "%"+params.Username+"%")
	}
	if params.Nickname != "" {
		db = db.Where("nickname like ? ", "%"+params.Nickname+"%")
	}

	err := db.Limit(params.PageSize).Offset(params.PageSize * (params.Page - 1)).Order("create_time DESC").Find(&lists).Error
	if err != nil {
		return err, nil, 0
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, nil, 0
	}

	return nil, lists, total
}
