package model

type Asset struct {
	Id                int     `json:"id" gorm:"column:id"`
	Name              string  `json:"name" gorm:"column:name"`                               //名称
	Remark            string  `json:"remark" gorm:"column:remark"`                           //描述
	Score             float64 `json:"score" gorm:"column:score"`                             //评分
	ViewLimit         int     `json:"view_limit" gorm:"column:view_limit"`                   //青少年观影限制 0  - 不限制 1 – 限制
	ViewCities        string  `json:"view_cities" gorm:"column:view_cities"`                 //地区限制 （json）
	ViewTags          string  `json:"view_tags" gorm:"column:view_tags"`                     //数据标签 （json）
	Type              int     `json:"type" gorm:"column:type"`                               //类型 1-视频 2-图片
	CoverUrl          string  `json:"cover_url" gorm:"column:cover_url"`                     //封面图地址
	MovieUrl          string  `json:"movie_url" gorm:"column:movie_url"`                     //视频地址
	PicUrls           string  `json:"pic_urls" gorm:"column:pic_urls"`                       //多张图片url  （json）
	DurationStartDate string  `json:"duration_start_date" gorm:"column:duration_start_date"` //资源有效期 - 开始时间
	DurationEndDate   string  `json:"duration_end_date" gorm:"column:duration_end_date"`     //资源有效期 - 结束时间
	ActType           int     `json:"act_type" gorm:"column:act_type"`                       //OK键动作类型 1-无动作 2-打开/下载应用 3-弹出二维码 4-加载长视频
	ActToast          string  `json:"act_toast" gorm:"column:act_toast"`                     //OK键引导文案
	ActOpenApps       string  `json:"act_open_apps" gorm:"column:act_open_apps"`             //需要下载打开的应用  （json）
	ActQrcodeUrl      string  `json:"act_qrcode_url" gorm:"column:act_qrcode_url"`           //二维码地址
	ActQrcodeOrgUrl   string  `json:"act_qrcode_org_url" gorm:"column:act_qrcode_org_url"`   //二维码原链接url
	ActQrcodeBgUrl    string  `json:"act_qrcode_bg_url" gorm:"column:act_qrcode_bg_url"`     //二维码背景图url
	ActPopTime        int     `json:"act_pop_time" gorm:"column:act_pop_time"`               //二维码自动弹出时间，单位：秒
	ActLongMovieUrl   string  `json:"act_long_movie_url" gorm:"column:act_long_movie_url"`   //长视频 url
	ShelfStatus       int     `json:"shelf_status" gorm:"column:shelf_status"`               //上架状态 1-未上架 2-已上架 3-已下架
	LastUpdateTime    string  `json:"last_update_time" gorm:"column:last_update_time"`
}
