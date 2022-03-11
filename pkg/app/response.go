package app

import (
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/crypt"
	"go-trailer-api/pkg/e"
	"go-trailer-api/pkg/model"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ListingResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}

func (g *Gin) ResponseEncryptJson(httpCode int, data []byte, key []byte) {
	msg, err := crypt.Packv2(&crypt.PData{
		Data: data,
		Key:  key,
	})
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ErrorEncryptError, nil)
		return
	}

	g.C.JSON(httpCode, model.EDataResponse{
		EK: string(msg.EK),
		ED: string(msg.ED),
		IV: string(msg.IV),
	})
}

func (g *Gin) ResponseJson(httpCode int, data interface{}) {
	if data == nil {
		data = CreateEmptyJson(httpCode)
	}

	g.C.JSON(httpCode, data)
}

type EmptyJson struct {
	Code int    `json:"code"`
	Msg  string `json:"mgs"`
}

func CreateEmptyJson(code int) *EmptyJson {
	return &EmptyJson{
		Code: code,
		Msg:  "",
	}
}
