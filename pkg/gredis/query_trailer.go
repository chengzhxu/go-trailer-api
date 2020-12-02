package gredis

type TrailerListParam struct {
	pageSize int `json:"page_size" binding:"" example:"20"`
	page int `json:"page" binding:"required" example:"1"`
	channelCode  string `json:"channel_code" binding:"required" example:""`
	deviceNo  string `json:"device_no" binding:"required" example:""`
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
