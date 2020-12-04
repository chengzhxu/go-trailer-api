package gredis

import (
	"encoding/json"
	"fmt"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/util"
	"strconv"
)

type Asset struct {
	Id                int     `json:"id" binding:"required,bas_primary"`
	Name              string  `json:"name" binding:"required"`                         //名称
	Remark            string  `json:"remark" binding:"required"`                       //描述
	Score             float64 `json:"score" binding:"required,min=0,max=10"`           //评分
	ViewLimit         int     `json:"view_limit" binding:""`                           //青少年观影限制 0  - 不限制 1 – 限制
	ViewCities        string  `json:"view_cities" binding:""`                          //地区限制 （json）
	ViewTags          string  `json:"view_tags" binding:""`                            //数据标签 （json）
	Type              int     `json:"type" binding:"required,asset_type"`              //类型 1-视频 2-图片
	CoverUrl          string  `json:"cover_url" binding:"required"`                    //封面图地址
	MovieUrl          string  `json:"movie_url" binding:"required"`                    //视频地址
	PicUrls           string  `json:"pic_urls" binding:"required"`                     //多张图片url  （json）
	DurationStartDate string  `json:"duration_start_date" binding:"required,bas_date"` //资源有效期 - 开始时间
	DurationEndDate   string  `json:"duration_end_date" binding:"required,bas_date"`   //资源有效期 - 结束时间
	ActType           int     `json:"act_type" binding:"required,asset_act_type"`      //OK键动作类型 1-无动作 2-打开/下载应用 3-弹出二维码 4-加载长视频
	ActToast          string  `json:"act_toast" binding:"required"`                    //OK键引导文案
	ActOpenApps       string  `json:"act_open_apps" binding:"required"`                //需要下载打开的应用  （json）
	ActQrcodeUrl      string  `json:"act_qrcode_url" binding:"required"`               //二维码地址
	ActQrcodeOrgUrl   string  `json:"act_qrcode_org_url" binding:"required"`           //二维码原链接url
	ActQrcodeBgUrl    string  `json:"act_qrcode_bg_url" binding:"required"`            //二维码背景图url
	ActPopTime        int     `json:"act_pop_time" binding:"required"`                 //二维码自动弹出时间，单位：秒
	ActLongMovieUrl   string  `json:"act_long_movie_url" binding:"required"`           //长视频 url
	ShelfStatus       int     `json:"shelf_status" binding:"required"`                 //上架状态 1-未上架 2-已上架 3-已下架
	LastUpdateTime    string  `json:"last_update_time" binding:"required,bas_time"`
}

type TrailerListParam struct {
	pageSize    int    `json:"page_size" binding:"" example:"20"`
	page        int    `json:"page" binding:"required" example:"1"`
	channelCode string `json:"channel_code" binding:"required" example:""`
	deviceNo    string `json:"device_no" binding:"required" example:""`
}

func mapAsset(a *Asset) map[string]interface{} {
	return map[string]interface{}{
		"id":                  a.Id,
		"name":                a.Name,
		"remark":              a.Remark,
		"score":               a.Score,
		"view_limit":          a.ViewLimit,
		"view_cities":         a.ViewCities,
		"view_tags":           a.ViewTags,
		"type":                a.Type,
		"cover_url":           a.CoverUrl,
		"pic_urls":            a.PicUrls,
		"duration_start_date": a.DurationStartDate,
		"duration_end_date":   a.DurationEndDate,
		"act_type":            a.ActType,
		"act_toast":           a.ActToast,
		"act_open_apps":       a.ActOpenApps,
		"act_qrcode_url":      a.ActQrcodeUrl,
		"act_qrcode_org_url":  a.ActQrcodeOrgUrl,
		"act_qrcode_bg_url":   a.ActQrcodeBgUrl,
		"act_pop_time":        a.ActPopTime,
		"act_long_movie_url":  a.ActLongMovieUrl,
		"shelf_status":        a.ShelfStatus,
		"last_update_time":    a.LastUpdateTime,
	}
}

/*
* 同步 Asset 数据并写入 Redis （暂定）
 */
func (asset *Asset) SyncAssetToRedis() error {
	lastUpdateTime := asset.LastUpdateTime
	tu := util.TimeToUnix(lastUpdateTime) //时间戳
	assetId := asset.Id
	asset.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", asset.Score), 64)
	jsonBytes, err := json.Marshal(asset)
	if err != nil {
		logging.Error(err)
		return err
	}

	conn := RedisConn.Get()
	_, err = conn.Do("zadd", "asset_id", assetId, assetId)
	if err != nil {
		logging.Error(err)
		return err
	}
	_, err = conn.Do("zadd", "asset_last_update_time", tu, assetId)
	if err != nil {
		logging.Error(err)
		return err
	}
	_, err = conn.Do("hset", "trailer_asset", assetId, jsonBytes)
	if err != nil {
		logging.Error(err)
		return err
	}

	return nil
}

/*
* 获取预告片
 */
//func QueryTotalAppUV(rr *TrailerListParam) (TrailerListParam, error) {
//	//var err error
//
//	pageSize := rr.pageSize
//	if pageSize == 0{
//		pageSize = 20
//	}
//	//page := rr.page
//
//
//
//	return nil, nil
//}
