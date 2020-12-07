package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-trailer-api/pkg/setting"
	"log"
)

var statsMySql *setting.StatsDatabase
var trailerMySql *setting.TrailerDatabase
var db *gorm.DB
var trailerDb *gorm.DB

// Setup initializes the database instance
func Setup() {
	statsMySql = setting.StatsDbSetting
	trailerMySql = setting.TrailerDbSetting
	createStatsConn()
	createTrailerConn()
}

// CreateStatsConn closes database connection (unnecessary)
func createStatsConn() {
	var err error
	db, err = gorm.Open(statsMySql.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		statsMySql.User,
		statsMySql.Password,
		statsMySql.Host,
		statsMySql.Name))
	if err != nil {
		log.Fatalf("Stats MySql.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.StatsDbSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

// CloseStatsDB closes database connection (unnecessary)
func CloseStatsDB() {
	defer db.Close()
}

// CreateTrailerConn closes database connection (unnecessary)
func createTrailerConn() {
	var err error
	trailerDb, err = gorm.Open(trailerMySql.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		trailerMySql.User,
		trailerMySql.Password,
		trailerMySql.Host,
		trailerMySql.Name))
	if err != nil {
		log.Fatalf("Trailer MySql.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(trailerDb *gorm.DB, defaultTableName string) string {
		return setting.TrailerDbSetting.TablePrefix + defaultTableName
	}

	trailerDb.SingularTable(true)
	//trailerDb.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//trailerDb.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//trailerDb.Callback().Delete().Replace("gorm:delete", deleteCallback)
	trailerDb.DB().SetMaxIdleConns(10)
	trailerDb.DB().SetMaxOpenConns(100)
}

// CloseDB closes database connection (unnecessary)
func CloseTrailerDB() {
	defer trailerDb.Close()
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
