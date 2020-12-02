package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-trailer-api/pkg/setting"
	"log"
)

var settingMySql *setting.Database
var db *gorm.DB

// Setup initializes the database instance
func Setup() {
	settingMySql = setting.DatabaseSetting
	createConn()
}

// CreateConn closes database connection (unnecessary)
func createConn() {
	var err error
	db, err = gorm.Open(settingMySql.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settingMySql.User,
		settingMySql.Password,
		settingMySql.Host,
		settingMySql.Name))

	if err != nil {
		log.Fatalf("MySql.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}



//func exec(sqlStr string) (*gorm.DB, error) {
//	response := db.Raw(sqlStr)
//
//	//defer func() {
//	//	db.Close()
//	//}()
//
//	return response, err
//}

