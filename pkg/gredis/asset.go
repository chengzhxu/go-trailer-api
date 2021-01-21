package gredis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/util"
	"math"
	"strconv"
	"strings"
)

type Asset struct {
	Id                 int         `json:"id" binding:"required,bas_primary"`
	Name               string      `json:"name" binding:""`                                   //名称
	Remark             string      `json:"remark" binding:""`                                 //描述
	Score              interface{} `json:"score" binding:""`                                  //评分
	ViewLimit          int         `json:"view_limit" binding:""`                             //青少年观影限制 0  - 不限制 1 – 限制
	ViewCities         string      `json:"view_cities" binding:"asset_region_code"`           //地区 （json）
	ViewTags           string      `json:"view_tags" binding:""`                              //数据标签 （json）
	Type               int         `json:"type" binding:"required,asset_type"`                //类型 1-视频 2-图片
	CoverUrl           string      `json:"cover_url" binding:""`                              //封面图地址
	MovieUrl           string      `json:"movie_url" binding:""`                              //视频地址
	PicUrls            interface{} `json:"pic_urls" binding:""`                               //多张图片url  （json）
	DurationStartDate  string      `json:"duration_start_date" binding:"required,bas_time"`   //资源有效期 - 开始时间
	DurationEndDate    string      `json:"duration_end_date" binding:"required,bas_time"`     //资源有效期 - 结束时间
	ActType            int         `json:"act_type" binding:"required,asset_act_type"`        //OK键动作类型 1-无动作 2-打开/下载应用 3-弹出二维码 4-加载长视频
	ActToast           string      `json:"act_toast" binding:""`                              //OK键引导文案
	Priority           int         `json:"priority" binding:""`                               //优先级  0:优先调用  1:优先下载
	ImgStayTime        int         `json:"img_stay_time" binding:""`                          //单张图片停留时长(秒)
	ActOpenApps        interface{} `json:"act_open_apps" binding:"asset_open_apps"`           //需要下载打开的应用  （json）
	ActQrcodeUrl       string      `json:"act_qrcode_url" binding:""`                         //二维码地址
	ActQrcodeOrgUrl    string      `json:"act_qrcode_org_url" binding:""`                     //二维码原链接url
	ActQrcodeBgUrl     string      `json:"act_qrcode_bg_url" binding:""`                      //二维码背景图url
	ActPopTime         int         `json:"act_pop_time" binding:""`                           //二维码自动弹出时间，单位：秒
	ActLongMovieUrl    string      `json:"act_long_movie_url" binding:""`                     //长视频 url
	ShelfStatus        int         `json:"shelf_status" binding:"required"`                   //上架状态 1-未上架 2-已上架 3-已下架
	LastUpdateTime     string      `json:"last_update_time" binding:"required,bas_time"`      //最后更新时间 - 排序使用
	DelFlag            int         `json:"del_flag" binding:"asset_is_del"`                   //是否删除  0:否  1:是
	ScreenControlCycle int         `json:"screen_control_cycle" binding:"" example:"0"`       //频控周期  天数  自然日
	ScreenControlCount int         `json:"screen_control_count" binding:"" example:"0"`       //周期内频控次数
	DisplayOrder       int         `json:"display_order" binding:"" example:"0"`              //排序
	ChannelCode        string      `json:"channel_code" binding:"asset_channel_code"`         //对应的渠道 - 全部为 ALL - json数组
	BanChannelCode     string      `json:"ban_channel_code" binding:"asset_ban_channel_code"` //排除渠道 - json数组
	OwnChannelId       int         `json:"own_channel_id" binding:""`                         //素材所属渠道
}

//预告片列表参数
type TrailerListParam struct {
	PageSize    int    `json:"page_size" binding:"" example:"20"`          //每页数量
	Page        int    `json:"page" binding:"required" example:"1"`        //页码
	ChannelCode string `json:"channel_code" binding:"required" example:""` //渠道码
	DeviceNo    string `json:"device_no" binding:"required" example:""`    //设备号
	IsSecure    bool   `json:"isSecure" binding:""`                        //判断返回链接形式 https or http;
	RegionCode  string `json:"region_code" binding:""`                     //region_code
}

type AssetArray []*Asset

type AssetResult struct {
	TotalPage   int        `json:"total_page" `
	TotalRows   int        `json:"total_rows" `
	CurrentPage int        `json:"current_page" `
	AssetArray  AssetArray `json:"list" `
}

//Redis key
var IdKey = "asset_id"                      //素材 id
var LastTimeKey = "asset_last_update_time"  //素材最后更新时间排序 - 已弃用
var DisplayOrderKey = "asset_display_order" //素材自定义排序 - 启用
var HSetKey = "trailer_asset"               //素材列表

var TrailerCapBeginKey = "trailer:"   //素材频控 redis key  first_name
var TrailerCapEndKey = ":asset:cap"   //素材频控 redis key  last_name
var TrailerCapCountEndKey = ":count:" //素材频控 - 周期内播放次数 redis key  last_name
var TrailerCapDateEndKey = ":date"    //素材频控 -  周期内首次播放时间 redis key  last_name

//定向 & 排除   Redis key
var TrailerGoalBeginKey = "trailer_asset:"             //素材定向目的 first_name
var ChannelGoal = TrailerGoalBeginKey + "channel:goal" //定向渠道
var ChannelBan = TrailerGoalBeginKey + "channel:ban"   //排除渠道
var RegionGoal = TrailerGoalBeginKey + "region:goal"   //定向地域

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
		"ban_channel_code":     a.BanChannelCode,
		"own_channel_id":       a.OwnChannelId,
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
	assetId := asset.Id
	conn := RedisConn.Get() //获取 Redis
	asset.DurationStartDate = util.CheckTime(asset.DurationStartDate)
	asset.DurationEndDate = util.CheckTime(asset.DurationEndDate)
	if !checkRunAsset(asset) { //无效的 Asset    -    从 Redis 删除
		err := delAssetForRedis(conn, assetId)
		if err != nil {
			logging.Error(err)
			return err
		}
	}

	delErr := delAssetReConf(conn, assetId) //删除 Asset 定向配置信息
	if delErr != nil {
		logging.Error(delErr)
		return delErr
	}

	err := setAssetForRedis(conn, asset) //写入 Asset
	if err != nil {
		logging.Error(err)
		return err
	}

	setErr := setAssetReConf(conn, asset) //写入 Asset 定向配置信息
	if setErr != nil {
		logging.Error(setErr)
		return setErr
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

	assetRes := AssetResult{}   //返回数据
	assetArr := AssetArray{}    //当前页素材信息
	capAssetArr := AssetArray{} //当素材全部频控后 - 重新返回的素材信息

	goalChannelArr := []string{} //定向 channel
	exChannelArr := []string{}   //排除 channel
	goalRegionArr := []string{}  //定向 region
	runAsset := []string{}       //定向 & 排除后的 asset

	conn := RedisConn.Get()
	//按照自定义倒序获取所有的 AssetId
	allAssetId, err := redis.Values(conn.Do("zrevrangebyscore", DisplayOrderKey, "+inf", "-inf"))
	if err != nil {
		logging.Error(err)
		return assetRes, err
	}
	goalChannelArr = getGoalChannel(conn, rr.ChannelCode) //定向 channel
	exChannelArr = getBanChannel(conn, rr.ChannelCode)    //排除 channel
	goalRegionArr = getGoalRegion(conn, rr.RegionCode)
	//定向 & 排除
	for _, v := range allAssetId {
		v := strings.ReplaceAll(fmt.Sprintf("%s", v), " ", "")
		if util.StrInArray(v, exChannelArr) { //排除 channel
			continue
		}
		if util.StrInArray(v, goalChannelArr) && util.StrInArray(v, goalRegionArr) { //定向 channel && 定向 region
			runAsset = append(runAsset, v)
		}
	}
	logging.Info(runAsset)
	onAsset := getOnAsset(conn, runAsset)

	rmCount := (page - 1) * pageSize //当前页之前的数量（忽略）
	totalRows := 0                   //总数量
	capTotalRows := 0                //频控均已到达时重新计算的总数量

	for _, v := range onAsset {
		v := fmt.Sprintf("%s", v)

		reAsset, _ := redis.Bytes(v, err)
		if err != nil {
			logging.Error(err)
			return assetRes, err
		}
		var asset *Asset
		json.Unmarshal(reAsset, &asset)

		if !checkRunAsset(asset) { //无效的 Asset
			continue
		}

		if !checkAssetChannel(asset, rr.ChannelCode) { //未在所属渠道
			continue
		}

		totalRows++      //总数量 - 排除 无效的数据
		capTotalRows++   //频控均已到达时重新计算的总数量
		if rmCount > 0 { //当前页之前的数据量 - 排除
			rmCount--
			continue
		}

		asset = tranceAsset(asset, rr.IsSecure)
		if !checkAssetCap(asset, rr.DeviceNo) { //到达频控 - 仅检查当前页及后的数据
			totalRows-- //总数量 - 排除 频控数据

			if len(capAssetArr) < pageSize { //素材全部频控时，重新返回的素材信息
				capAssetArr = append(capAssetArr, asset)
			}
			continue
		}

		if len(assetArr) < pageSize { //当前页数据
			assetArr = append(assetArr, asset)
		}
	}

	if len(assetArr) == 0 { //所有素材均已到达频控 -  重新返回
		assetArr = capAssetArr   //重新获取的素材信息
		totalRows = capTotalRows //重新计算页码
	}
	pageCount := int(math.Ceil(float64(totalRows) / float64(pageSize))) //总页数

	assetRes.TotalPage = pageCount
	assetRes.TotalRows = totalRows
	assetRes.CurrentPage = page
	assetRes.AssetArray = assetArr

	return assetRes, nil
}

//封装客户端返回的素材信息
func tranceAsset(asset *Asset, isSecure bool) *Asset {
	var mapPics []interface{} //图片集合
	if asset.PicUrls != "" && asset.PicUrls != nil {
		err := json.Unmarshal([]byte(asset.PicUrls.(string)), &mapPics)
		if err != nil {
			logging.Error(err)
			mapPics = make([]interface{}, 0)
		} else {
			if !isSecure { //替换 https =》 http
				mapPics = replaceAssetPicUrlsHttp(mapPics)
			}
		}
	}
	asset.PicUrls = mapPics

	oaArr := []OpenApp{} //打开的app
	if asset.ActOpenApps != "" && asset.ActOpenApps != nil {
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

	if !isSecure { //替换 https =》 http
		asset = replaceAssetHttp(asset)
	}

	return asset
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

//清除 Asset    与 setAssetForRedis 方法中的 redis key 保持一致
func delAssetForRedis(conn redis.Conn, assetId int) error {
	_, err := conn.Do("zrem", IdKey, assetId) //从 Id 排序中移除
	if err != nil {
		return err
	}
	_, oe := conn.Do("zrem", DisplayOrderKey, assetId) //从自定义排序中移除
	if oe != nil {
		return oe
	}
	_, he := conn.Do("hdel", HSetKey, assetId) //从 Hash 中移除
	if he != nil {
		return he
	}

	return nil
}

//写入 Asset    与 delAssetForRedis 方法中的 redis key 保持一致
func setAssetForRedis(conn redis.Conn, asset *Asset) error {
	assetId := asset.Id
	asset.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", asset.Score), 64)
	jsonBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	_, err = conn.Do("zadd", IdKey, assetId, assetId) //根据 Id 排序
	if err != nil {
		return err
	}
	_, err = conn.Do("zadd", DisplayOrderKey, asset.DisplayOrder, assetId) //根据自定义排序
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", HSetKey, assetId, jsonBytes) //根据 Id 存储
	if err != nil {
		return err
	}

	return nil
}

// 素材 url 头替换  https =》 http
func replaceAssetHttp(asset *Asset) *Asset {
	domainArr := []string{".shafa.com", ".xmxapi.com", ".xmxcdn.com"} //自己域名不替换

	hs := "https://"
	h := "http://"

	if strings.LastIndex(asset.MovieUrl, hs) > -1 { //视频url
		flag := false
		for _, v := range domainArr {
			if strings.LastIndex(asset.MovieUrl, v) > -1 { //   自己服务器资源
				flag = true
				break
			}
		}
		if flag {
			asset.MovieUrl = strings.Replace(asset.MovieUrl, hs, h, 1)
		}
	}
	if strings.LastIndex(asset.CoverUrl, hs) > -1 { // 缩略图
		flag := false
		for _, v := range domainArr {
			if strings.LastIndex(asset.CoverUrl, v) > -1 { //   自己服务器资源
				flag = true
				break
			}
		}
		if flag {
			asset.CoverUrl = strings.Replace(asset.CoverUrl, hs, h, 1)
		}
	}
	if strings.LastIndex(asset.ActQrcodeUrl, hs) > -1 { //二维码 url
		flag := false
		for _, v := range domainArr {
			if strings.LastIndex(asset.ActQrcodeUrl, v) > -1 { //   自己服务器资源
				flag = true
				break
			}
		}
		if flag {
			asset.ActQrcodeUrl = strings.Replace(asset.ActQrcodeUrl, hs, h, 1)
		}
	}
	if strings.LastIndex(asset.ActQrcodeOrgUrl, hs) > -1 { //二维码原 url
		flag := false
		for _, v := range domainArr {
			if strings.LastIndex(asset.ActQrcodeOrgUrl, v) > -1 { //   自己服务器资源
				flag = true
				break
			}
		}
		if flag {
			asset.ActQrcodeOrgUrl = strings.Replace(asset.ActQrcodeOrgUrl, hs, h, 1)
		}
	}
	if strings.LastIndex(asset.ActQrcodeBgUrl, hs) > -1 { //二维码背景 url
		flag := false
		for _, v := range domainArr {
			if strings.LastIndex(asset.ActQrcodeBgUrl, v) > -1 { //   自己服务器资源
				flag = true
				break
			}
		}
		if flag {
			asset.ActQrcodeBgUrl = strings.Replace(asset.ActQrcodeBgUrl, hs, h, 1)
		}
	}

	if strings.LastIndex(asset.ActLongMovieUrl, hs) > -1 { //长视频url
		flag := false
		for _, v := range domainArr {
			if strings.LastIndex(asset.ActLongMovieUrl, v) > -1 { //   自己服务器资源
				flag = true
				break
			}
		}
		if flag {
			asset.ActLongMovieUrl = strings.Replace(asset.ActLongMovieUrl, hs, h, 1)
		}
	}

	return asset
}

func replaceAssetPicUrlsHttp(mapPics []interface{}) []interface{} {
	if len(mapPics) > 0 {
		domainArr := []string{".shafa.com", ".xmxapi.com", ".xmxcdn.com"} //自己域名不替换
		hs := "https://"
		h := "http://"

		for i, v := range mapPics {
			pu := v.(string)
			if strings.LastIndex(pu, hs) > -1 {
				flag := false
				for _, a := range domainArr {
					if strings.LastIndex(pu, a) > -1 { //   自己服务器资源
						flag = true
						break
					}
				}
				if flag {
					pu = strings.Replace(pu, hs, h, 1)
					mapPics[i] = pu
				}
			}
		}
	}

	return mapPics
}

//删除素材定向 or 排除信息
func delAssetReConf(conn redis.Conn, assetId int) error {
	if assetId > 0 {
		reply, err := conn.Do("hget", HSetKey, assetId)
		if reply == nil {
			return nil
		}
		reAsset, _ := redis.Bytes(reply, err)
		if err != nil {
			logging.Error(err)
			return err
		}

		var asset *Asset
		json.Unmarshal(reAsset, &asset) //之前存放在 redis 中的素材信息
		e := todoAssetReConf(asset, 0, conn)
		if e != nil {
			logging.Error(e)
			return e
		}
	}

	return nil
}

//写入素材定向 or 排除信息
func setAssetReConf(conn redis.Conn, asset *Asset) error {
	err := todoAssetReConf(asset, 1, conn)
	if err != nil {
		logging.Error(err)
		return err
	}

	return nil
}

//素材定向 or 排除配置  redis zset
func todoAssetReConf(asset *Asset, t int, conn redis.Conn) error {
	if asset.Id > 0 {
		assetIdStr := strconv.Itoa(asset.Id)
		redisVal := ":" + assetIdStr //redis value 固定部分
		redisScore := 0              //分值 0 - 建议统一分值

		if asset.ChannelCode != "" { //定向 channel
			var mapChannel []string
			err := json.Unmarshal([]byte(asset.ChannelCode), &mapChannel)
			if err != nil {
				logging.Error(err)
				return err
			}
			for _, c := range mapChannel {
				v := c + redisVal //redis value 规则
				if t == 0 {       //删除
					_, err := conn.Do("zrem", ChannelGoal, v)
					if err != nil {
						logging.Error(err)
						return err
					}
				} else if t == 1 { //新增
					_, err := conn.Do("zadd", ChannelGoal, redisScore, v)
					if err != nil {
						logging.Error(err)
						return err
					}
				}

			}
		}

		if asset.BanChannelCode != "" { //排除 channel
			var mapChannel []string
			err := json.Unmarshal([]byte(asset.BanChannelCode), &mapChannel)
			if err != nil {
				logging.Error(err)
				return err
			}
			for _, c := range mapChannel {
				v := c + redisVal //redis value 规则
				if t == 0 {       //删除
					_, err := conn.Do("zrem", ChannelBan, v)
					if err != nil {
						logging.Error(err)
						return err
					}
				} else if t == 1 { //新增
					_, err := conn.Do("zadd", ChannelBan, redisScore, v)
					if err != nil {
						logging.Error(err)
						return err
					}
				}
			}
		}

		if asset.ViewCities != "" { //定向地区
			var mapRegion []string
			err := json.Unmarshal([]byte(asset.ViewCities), &mapRegion)
			if err != nil {
				logging.Error(err)
				return err
			}
			for _, r := range mapRegion {
				v := r + redisVal //redis value 规则
				if t == 0 {       //删除
					_, err := conn.Do("zrem", RegionGoal, v)
					if err != nil {
						logging.Error(err)
						return err
					}
				} else if t == 1 { //新增
					_, err := conn.Do("zadd", RegionGoal, redisScore, v)
					if err != nil {
						logging.Error(err)
						return err
					}
				}
			}
		}
	}

	return nil
}

//定向 channel
func getGoalChannel(conn redis.Conn, channelCode string) []string {
	assetArr := []string{}

	if channelCode != "" {
		channelStr := "(" + channelCode + ":"
		goalChannel, _ := redis.Values(conn.Do("ZRANGEBYLEX", ChannelGoal, channelStr, channelStr+"\\xff")) //定向 channel

		for _, r := range goalChannel {
			v := fmt.Sprintf("%s", r)
			arr := strings.Split(v, ":")
			if len(arr) == 2 {
				assetArr = append(assetArr, strings.ReplaceAll(arr[1], " ", ""))
			}
		}
	}

	allChannel, _ := redis.Values(conn.Do("ZRANGEBYLEX", ChannelGoal, "(ALL:", "(ALL:\\xff")) //ALl channel
	for _, r := range allChannel {
		v := fmt.Sprintf("%s", r)
		arr := strings.Split(v, ":")
		if len(arr) == 2 {
			assetArr = append(assetArr, strings.ReplaceAll(arr[1], " ", ""))
		}
	}

	return assetArr
}

//排除 channel
func getBanChannel(conn redis.Conn, channelCode string) []string {
	assetArr := []string{}

	if channelCode != "" {
		channelStr := "(" + channelCode + ":"
		exChannel, _ := redis.Values(conn.Do("ZRANGEBYLEX", ChannelBan, channelStr, channelStr+"\\xff")) //排除 channel

		for _, r := range exChannel {
			v := fmt.Sprintf("%s", r)
			arr := strings.Split(v, ":")
			if len(arr) == 2 {
				assetArr = append(assetArr, strings.ReplaceAll(arr[1], " ", ""))
			}
		}
	}

	return assetArr
}

//定向 region
func getGoalRegion(conn redis.Conn, regionCode string) []string {
	assetArr := []string{}

	if regionCode != "" {
		regionStr := "(" + regionCode + ":"
		goalRegion, _ := redis.Values(conn.Do("ZRANGEBYLEX", RegionGoal, regionStr, regionStr+"\\xff")) //定向 region

		for _, r := range goalRegion {
			v := fmt.Sprintf("%s", r)
			arr := strings.Split(v, ":")
			if len(arr) == 2 {
				assetArr = append(assetArr, strings.ReplaceAll(arr[1], " ", ""))
			}
		}
	}

	allRegion, _ := redis.Values(conn.Do("ZRANGEBYLEX", RegionGoal, "(ALL:", "(ALL:\\xff")) //排除 region
	for _, r := range allRegion {
		v := fmt.Sprintf("%s", r)
		arr := strings.Split(v, ":")
		if len(arr) == 2 {
			assetArr = append(assetArr, strings.ReplaceAll(arr[1], " ", ""))
		}
	}

	return assetArr
}

//获取 Asset  -  根据定向及排除过滤后的
func getOnAsset(conn redis.Conn, items []string) []interface{} {
	if len(items) > 0 {
		fields := make([]interface{}, len(items)+1)
		fields[0] = HSetKey
		for i, vv := range items {
			fields[i+1] = vv
		}
		reply, err := conn.Do("hmget", fields...)

		if reply != nil {
			assetRes, _ := redis.Values(reply, err)
			return assetRes
		}

	}

	return nil
}

//清洗 Redis Asset
func ResetAsset() error {
	conn := RedisConn.Get() //获取 Redis

	assetArr, err := redis.Values(conn.Do("hkeys", HSetKey))
	if err != nil {
		logging.Error(err)
		return err
	}

	for _, v := range assetArr {
		assetId := fmt.Sprintf("%s", v)
		score, _ := conn.Do("ZSCORE", DisplayOrderKey, assetId) //获取自定义排序数据

		if score == nil { // 不存在自定义排序中  移除掉
			_, err := conn.Do("zrem", IdKey, assetId) //从 Id 排序中移除
			if err != nil {
				logging.Error(err)
				return err
			}

			_, he := conn.Do("hdel", HSetKey, assetId) //从 Hash 中移除
			if he != nil {
				logging.Error(he)
				return he
			}
		}
	}

	return nil
}

//移除 Redis Asset
func RemoveAsset(assetId int) error {
	conn := RedisConn.Get() //获取 Redis

	e := delAssetReConf(conn, assetId) //先 移除 Redis 中 Asset 的定向配置信息; 因为移除配置信息需要获取 Redis 中素材的配置信息
	if e != nil {
		logging.Error(e)
		return e
	}

	err := delAssetForRedis(conn, assetId) //后 移除 Redis 中的 Asset 信息
	if err != nil {
		logging.Error(err)
		return err
	}

	return nil
}
