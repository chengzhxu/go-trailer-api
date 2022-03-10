package model

import "github.com/jinzhu/gorm"

type AppBlackWhite struct {
	Id           int    `json:"id" gorm:"column:id"`
	AppId        int    `json:"app_id" gorm:"column:app_id"`                 //app_id
	AppName      string `json:"app_name" gorm:"column:app_name"`             //app_name
	BlackList    string `json:"black_list" gorm:"column:black_list"`         //黑名单规则
	WhiteList    string `json:"white_list" gorm:"column:white_list"`         //白名单规则
	CreateTime   string `json:"create_time" gorm:"column:create_time"`       //创建时间
	CreateUserId int    `json:"create_user_id" gorm:"column:create_user_id"` //创建人
	UpdateTime   string `json:"update_time" gorm:"column:update_time"`       //最后修改时间
	UpdateUserId int    `json:"update_user_id" gorm:"column:update_user_id"` //最后修改人
}

func (AppBlackWhite) TableName() string {
	return "app_version_blackwhite"
}

func GetBlackWhiteByAppId(appId int) ([]*AppBlackWhite, error) {
	var list []*AppBlackWhite
	if appId > 0 {
		err := trailerDb.Where("app_id = ?", appId).Find(&list).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return list, nil
}
