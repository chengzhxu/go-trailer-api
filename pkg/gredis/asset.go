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
	Id                 int         `json:"id" binding:"required,bas_primary"`
	Name               string      `json:"name" binding:""`                                 //名称
	Remark             string      `json:"remark" binding:""`                               //描述
	Score              interface{} `json:"score" binding:""`                                //评分
	ViewLimit          int         `json:"view_limit" binding:""`                           //青少年观影限制 0  - 不限制 1 – 限制
	ViewCities         string      `json:"view_cities" binding:""`                          //地区限制 （json）
	ViewTags           string      `json:"view_tags" binding:""`                            //数据标签 （json）
	Type               int         `json:"type" binding:"required,asset_type"`              //类型 1-视频 2-图片
	CoverUrl           string      `json:"cover_url" binding:""`                            //封面图地址
	MovieUrl           string      `json:"movie_url" binding:""`                            //视频地址
	PicUrls            interface{} `json:"pic_urls" binding:""`                             //多张图片url  （json）
	DurationStartDate  string      `json:"duration_start_date" binding:"required,bas_time"` //资源有效期 - 开始时间
	DurationEndDate    string      `json:"duration_end_date" binding:"required,bas_time"`   //资源有效期 - 结束时间
	ActType            int         `json:"act_type" binding:"required,asset_act_type"`      //OK键动作类型 1-无动作 2-打开/下载应用 3-弹出二维码 4-加载长视频
	ActToast           string      `json:"act_toast" binding:""`                            //OK键引导文案
	Priority           int         `json:"priority" binding:""`                             //优先级  0:优先调用  1:优先下载
	ImgStayTime        int         `json:"img_stay_time" binding:""`                        //单张图片停留时长(秒)
	ActOpenApps        interface{} `json:"act_open_apps" binding:"asset_open_apps"`         //需要下载打开的应用  （json）
	ActQrcodeUrl       string      `json:"act_qrcode_url" binding:""`                       //二维码地址
	ActQrcodeOrgUrl    string      `json:"act_qrcode_org_url" binding:""`                   //二维码原链接url
	ActQrcodeBgUrl     string      `json:"act_qrcode_bg_url" binding:""`                    //二维码背景图url
	ActPopTime         int         `json:"act_pop_time" binding:""`                         //二维码自动弹出时间，单位：秒
	ActLongMovieUrl    string      `json:"act_long_movie_url" binding:""`                   //长视频 url
	ShelfStatus        int         `json:"shelf_status" binding:"required"`                 //上架状态 1-未上架 2-已上架 3-已下架
	LastUpdateTime     string      `json:"last_update_time" binding:"required,bas_time"`    //最后更新时间 - 排序使用
	DelFlag            int         `json:"del_flag" binding:"asset_is_del"`                 //是否删除  0:否  1:是
	ScreenControlCycle int         `json:"screen_control_cycle" binding:"" example:"0"`     //频控周期  天数  自然日
	ScreenControlCount int         `json:"screen_control_count" binding:"" example:"0"`     //周期内频控次数
	DisplayOrder       int         `json:"display_order" binding:"" example:"0"`            //排序
	ChannelCode        string      `json:"channel_code" binding:"asset_channel_code"`       //对应的渠道 - 全部为 ALL - json数组
}

//预告片列表参数
type TrailerListParam struct {
	PageSize    int    `json:"page_size" binding:"" example:"20"`          //每页数量
	Page        int    `json:"page" binding:"required" example:"1"`        //页码
	ChannelCode string `json:"channel_code" binding:"required" example:""` //渠道码
	DeviceNo    string `json:"device_no" binding:"required" example:""`    //设备号
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
var DisplayOrderKey = "asset_display_order"
var HSetKey = "trailer_asset"

var TrailerCapBeginKey = "trailer:"   //素材频控 redis key  first_name
var TrailerCapEndKey = ":asset:cap"   //素材频控 redis key  last_name
var TrailerCapCountEndKey = ":count:" //素材频控 - 周期内播放次数 redis key  last_name
var TrailerCapDateEndKey = ":date"    //素材频控 -  周期内首次播放时间 redis key  last_name

func mapAsset(a *Asset) map[string]interface{} {
	return map[string]interface{}{
		"id":                   a.Id,
		"name":                 a.Name,
		"remark":               a.Remark,
		"score":                a.Score,
		"view_limit":           a.ViewLimit,
		"view_cities":          a.ViewCities,
		"view_tags":            a.ViewTags,
		"type":                 a.Type,
		"cover_url":            a.CoverUrl,
		"pic_urls":             a.PicUrls,
		"duration_start_date":  a.DurationStartDate,
		"duration_end_date":    a.DurationEndDate,
		"act_type":             a.ActType,
		"act_toast":            a.ActToast,
		"priority":             a.Priority,
		"img_stay_time":        a.ImgStayTime,
		"act_open_apps":        a.ActOpenApps,
		"act_qrcode_url":       a.ActQrcodeUrl,
		"act_qrcode_org_url":   a.ActQrcodeOrgUrl,
		"act_qrcode_bg_url":    a.ActQrcodeBgUrl,
		"act_pop_time":         a.ActPopTime,
		"act_long_movie_url":   a.ActLongMovieUrl,
		"shelf_status":         a.ShelfStatus,
		"last_update_time":     a.LastUpdateTime,
		"del_flag":             a.DelFlag,
		"screen_control_cycle": a.ScreenControlCycle,
		"screen_control_count": a.ScreenControlCount,
		"display_order":        a.DisplayOrder,
		"channel_code":         a.ChannelCode,
	}
}

type OpenApp struct {
	App         string      `json:"app" `
	AppId       string      `json:"appId" `
	Type        string      `json:"type" `
	Schema      string      `json:"schema" `
	PackageName interface{} `json:"packageName" `
}

/*
* 同步 Asset 数据并写入 Redis （暂定）
 */
func (asset *Asset) SyncAssetToRedis() error {
	conn := RedisConn.Get() //获取 Redis
	asset.DurationStartDate = util.CheckTime(asset.DurationStartDate)
	asset.DurationEndDate = util.CheckTime(asset.DurationEndDate)
	if !checkRunAsset(asset) { //无效的 Asset    -    从 Redis 删除
		_, err := conn.Do("zrem", IdKey, asset.Id) //从 Id 排序中移除
		if err != nil {
			logging.Error(err)
			return err
		}
		//_, e := conn.Do("zrem", LastTimeKey, asset.Id) //从最后更新时间排序中移除
		//if e != nil {
		//	logging.Error(e)
		//	return e
		//}
		_, oe := conn.Do("zrem", DisplayOrderKey, asset.Id) //从自定义排序中移除
		if oe != nil {
			logging.Error(oe)
			return oe
		}
		_, he := conn.Do("hdel", HSetKey, asset.Id) //从 Hash 中移除
		if he != nil {
			logging.Error(he)
			return he
		}

		return nil
	}

	//lastUpdateTime := asset.LastUpdateTime
	//tu := util.TimeToUnix(lastUpdateTime) //时间戳
	assetId := asset.Id
	asset.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", asset.Score), 64)
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
	//_, err = conn.Do("zadd", LastTimeKey, tu, assetId) //根据最后更新时间排序
	//if err != nil {
	//	logging.Error(err)
	//	return err
	//}
	_, err = conn.Do("zadd", DisplayOrderKey, asset.DisplayOrder, assetId) //根据自定义排序
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
	//按照自定义倒序获取所有的 AssetId
	res, err := redis.Values(conn.Do("zrevrangebyscore", DisplayOrderKey, "+inf", "-inf"))
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

		if !checkAssetChannel(asset, rr.ChannelCode) { //未在所属渠道
			continue
		}

		totalRows++      //总数量 - 排除 无效的数据
		if rmCount > 0 { //当前页之前的数据量 - 排除
			rmCount--
			continue
		}

		if !checkAssetCap(asset, rr.DeviceNo) { //到达频控 - 仅检查当前页及后的数据
			totalRows-- //总数量 - 排除 频控数据
			continue
		}

		if len(assetArr) < pageSize { //当前页数据
			var mapPics []interface{} //图片集合
			if asset.PicUrls != "" {
				err := json.Unmarshal([]byte(asset.PicUrls.(string)), &mapPics)
				if err != nil {
					logging.Error(err)
					mapPics = make([]interface{}, 0)
				}
			}
			asset.PicUrls = mapPics

			oaArr := []OpenApp{} //打开的app
			if asset.ActOpenApps != "" {
				var openAppArray []OpenApp
				err := json.Unmarshal([]byte(asset.ActOpenApps.(string)), &openAppArray)
				if err != nil {
					logging.Error(err)
				} else {
					for _, p := range openAppArray {
						oaArr = append(oaArr, p)
					}
				}
			}
			asset.ActOpenApps = oaArr

			asset.Score = fmt.Sprintf("%0.1f", asset.Score) //评分 - 保留1位小数 若为整数，则为 x.0

			assetArr = append(assetArr, asset)
		}
	}

	pageCount := int(math.Ceil(float64(totalRows) / float64(pageSize))) //总页数

	assetRes.TotalPage = pageCount
	assetRes.TotalRows = totalRows
	assetRes.CurrentPage = page
	assetRes.AssetArray = assetArr

	return assetRes, nil
}

//检查 Asset 是否有效
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
	if asset.DelFlag == 1 { //删除状态
		return false
	}

	return true
}

//检查 Asset 所属渠道
func checkAssetChannel(asset *Asset, channelCode string) bool {
	if asset.ChannelCode != "" && channelCode != "" { //匹配 channel_code
		var mapChannel []string
		err := json.Unmarshal([]byte(asset.ChannelCode), &mapChannel)
		if err != nil {
			logging.Error(err)
			return false
		}
		if len(mapChannel) == 1 && mapChannel[0] == "ALL" { //全渠道
			return true
		}

		if !util.StrInArray(channelCode, mapChannel) { //未匹配到对应的 channel_code
			return false
		}
	}

	return true
}

//检查 Asset 频控
func checkAssetCap(asset *Asset, deviceNo string) bool {
	if asset.Id > 0 && asset.ScreenControlCount > 0 && deviceNo != "" {
		conn := RedisConn.Get()                                                            //获取 Redis
		TrailerCapKey := TrailerCapBeginKey + deviceNo + TrailerCapEndKey                  //频控 Redis key
		CapCountKey := TrailerCapBeginKey + strconv.Itoa(asset.Id) + TrailerCapCountEndKey //素材播放次数 Redis key
		count, _ := redis.Int(conn.Do("hget", TrailerCapKey, CapCountKey))
		if count >= asset.ScreenControlCount { //判断是否达到频控次数   播放次数 >= 频控
			return false
		}
	}

	return true
}
