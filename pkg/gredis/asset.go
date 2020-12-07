package gredis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/util"
	"math"
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
	DurationStartDate string  `json:"duration_start_date" binding:"required,bas_time"` //资源有效期 - 开始时间
	DurationEndDate   string  `json:"duration_end_date" binding:"required,bas_time"`   //资源有效期 - 结束时间
	ActType           int     `json:"act_type" binding:"required,asset_act_type"`      //OK键动作类型 1-无动作 2-打开/下载应用 3-弹出二维码 4-加载长视频
	ActToast          string  `json:"act_toast" binding:"required"`                    //OK键引导文案
	ActOpenApps       string  `json:"act_open_apps" binding:"required"`                //需要下载打开的应用  （json）
	ActQrcodeUrl      string  `json:"act_qrcode_url" binding:"required"`               //二维码地址
	ActQrcodeOrgUrl   string  `json:"act_qrcode_org_url" binding:"required"`           //二维码原链接url
	ActQrcodeBgUrl    string  `json:"act_qrcode_bg_url" binding:"required"`            //二维码背景图url
	ActPopTime        int     `json:"act_pop_time" binding:"required"`                 //二维码自动弹出时间，单位：秒
	ActLongMovieUrl   string  `json:"act_long_movie_url" binding:"required"`           //长视频 url
	ShelfStatus       int     `json:"shelf_status" binding:"required"`                 //上架状态 1-未上架 2-已上架 3-已下架
	LastUpdateTime    string  `json:"last_update_time" binding:"required,bas_time"`    //最后更新时间 - 排序使用
	IsDel             int     `json:"is_del" binding:"asset_is_del"`
}

//预告片列表参数
type TrailerListParam struct {
	PageSize    int    `json:"page_size" binding:"" example:"20"`
	Page        int    `json:"page" binding:"required" example:"1"`
	ChannelCode string `json:"channel_code" binding:"required" example:""`
	DeviceNo    string `json:"device_no" binding:"required" example:""`
}

type AssetArray []*Asset

type AssetResult struct {
	TotalPage   int        `json:"total_page" `
	TotalRows   int        `json:"total_rows" `
	CurrentPage int        `json:"current_page" `
	AssetArray  AssetArray `json:"list" `
}

var IdKey = "asset_id"
var LastTimeKey = "asset_last_update_time"
var HSetKey = "trailer_asset"

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
		"is_del":              a.IsDel,
	}
}

/*
* 同步 Asset 数据并写入 Redis （暂定）
 */
func (asset *Asset) SyncAssetToRedis() error {
	conn := RedisConn.Get()    //获取 Redis
	if !checkRunAsset(asset) { //无效的 Asset    -    从 Redis 删除
		_, err := conn.Do("zrem", IdKey, asset.Id) //从 Id 排序中移除
		if err != nil {
			logging.Error(err)
			return err
		}
		_, e := conn.Do("zrem", LastTimeKey, asset.Id) //从最后更新时间排序中移除
		if e != nil {
			logging.Error(e)
			return e
		}
		_, he := conn.Do("hdel", HSetKey, asset.Id) //从 Hash 中移除
		if he != nil {
			logging.Error(he)
			return he
		}

		return nil
	}

	lastUpdateTime := asset.LastUpdateTime
	tu := util.TimeToUnix(lastUpdateTime) //时间戳
	assetId := asset.Id
	asset.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", asset.Score), 64)
	jsonBytes, err := json.Marshal(asset)
	if err != nil {
		logging.Error(err)
		return err
	}

	_, err = conn.Do("zadd", IdKey, assetId, assetId) //根据 Id 排序
	if err != nil {
		logging.Error(err)
		return err
	}
	_, err = conn.Do("zadd", LastTimeKey, tu, assetId) //根据最后更新时间排序
	if err != nil {
		logging.Error(err)
		return err
	}
	_, err = conn.Do("hset", HSetKey, assetId, jsonBytes) //根据 Id 存储
	if err != nil {
		logging.Error(err)
		return err
	}

	return nil
}

/*
* 获取预告片
 */
func (rr *TrailerListParam) QueryTrailerList() (AssetResult, error) {
	var err error

	pageSize := rr.PageSize //每页数量
	if pageSize == 0 {
		pageSize = 20
	}
	page := rr.Page //页码

	assetRes := AssetResult{} //返回数据
	assetArr := AssetArray{}  //当前页素材信息

	conn := RedisConn.Get()
	//按照最后修改时间倒序获取所有的 AssetId
	res, err := redis.Values(conn.Do("zrevrangebyscore", LastTimeKey, 7889155200, 1))
	if err != nil {
		logging.Error(err)
		return assetRes, err
	}

	rmCount := (page - 1) * pageSize //当前页之前的数量（忽略）
	totalRows := 0                   //总数量

	for _, v := range res {
		if err != nil {
			logging.Error(err)
			return assetRes, err
		}

		reply, err := redis.Bytes(conn.Do("hget", HSetKey, v))
		if err != nil {
			logging.Error(err)
			return assetRes, err
		}

		var asset *Asset
		json.Unmarshal(reply, &asset)

		if !checkRunAsset(asset) { //无效的 Asset
			continue
		}

		if rmCount > 0 { //当前页之前的数据量
			rmCount--
			continue
		}
		if len(assetArr) < pageSize {
			assetArr = append(assetArr, asset)
		}
		totalRows++
	}

	pageCount := int(math.Ceil(float64(totalRows) / float64(pageSize))) //总页数

	assetRes.TotalPage = pageCount
	assetRes.TotalRows = totalRows
	assetRes.CurrentPage = page
	assetRes.AssetArray = assetArr

	return assetRes, nil
}

//判断 Asset 是否有效
func checkRunAsset(asset *Asset) bool {
	if asset.ShelfStatus != 2 { // 未上架状态
		return false
	}
	start_time := util.TimeToUnix(asset.DurationStartDate) //资源有效开始时间
	end_time := util.TimeToUnix(asset.DurationEndDate)     //资源有效结束时间
	new_time := util.GetNowTimeStamp()
	if new_time < start_time || new_time > end_time { //不在有效期时间内
		return false
	}
	if asset.IsDel == 0 { //删除状态
		return false
	}

	return true
}
