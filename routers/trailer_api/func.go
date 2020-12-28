package trailer_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-trailer-api/pkg/app"
	"go-trailer-api/pkg/crypt"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/model"
	"net/http"
)

var privateKeyBytes = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCyBUxPA9X3tdJSsgNi0/CctBsBbvJEUgt5FWbbG4nTp2r7XK2T
Vk/YBZ6jq7kiA1rzJ/xmJTdpUdVUhPfX1DN/7iapb5K3z/NhIO/jsJgGO+YtgK5I
WEcjGwElPZnOsMk6iNZAWCGa7EyA0FHkul7w4eFOjC+RGqbKfsl306EUtwIDAQAB
AoGBAIMp8ip5ugoUVk4FyQbk/3CGJyusMiZyiO+C/FDN/oQK44EmrOFVA+k3YsZW
/UX5UOa9fHNKUoRv/g2TFwVX3UTVuri4h4x6zPfCOv3O1TUVaJzgCGa8bc/sNRu4
AlYZZdeKNxQbURbht/kbyu11HjyPqPP+ZXcwoalCiZQmtI/xAkEA7VbfGeQuH8Go
4SBv3EbQvFhJBWjjxmjIq3VQZD7HqdZMTqYnqaq3bVvqHWizrxyku2BbgaFbV5MJ
ZhbeZOF/bwJBAMAEeGbjs9xkPcfMcOCVypZpRcoi/0BMzj6Yl8YKmzxY6cNVIyEB
UukklSxZlbt8piFaZHZo1e85g1pDiX01GzkCQBvTx6y9eDr49dgPeY4WL3slzsn3
ll05A+42fwqB4d8j5SaDjLrz7TXBRR3VnNu3PAlMLu5wAMmvz7ZMkB674bkCQQCE
ZlCy+Uz6qW/kBXbLlN2Eyv/hOjJwnsUTalo0pvmVKeW910WKq4QE2EG3u+m/xloy
40YkU3M4KZsFsU3rNKQZAkAdJJEhk3+LwoCSD/hneL7aCtCAmmmWOuHTWjQAD1Zo
OEJTQoc9/wQ767QpGcbqdwF+KjaEVTqn5Wx+9UkMWhyy
-----END RSA PRIVATE KEY-----
`)

func BindParameter(c *gin.Context) (*model.EData, error) {
	enReq := &model.EData{}
	err := c.ShouldBindJSON(enReq)
	if err != nil {
		logging.Info(err)
		return nil, err
	}
	return enReq, err
}

func GinDecryptData(c *gin.Context, appG app.Gin) (*crypt.PData, error) {
	enReq, err := BindParameter(c)
	if err != nil {
		appG.ResponseJson(http.StatusBadRequest, err)
		return nil, err
	}

	pData, err := Unpack(enReq)
	if err != nil {
		appG.ResponseJson(http.StatusBadRequest, err)
		return nil, err
	}
	return pData, err
}

func ParamDecryptData(enReq *model.EData, appG app.Gin) (*crypt.PData, error) {
	pData, err := Unpack(enReq)
	if err != nil {
		appG.ResponseJson(http.StatusBadRequest, nil)
		return nil, err
	}
	return pData, err
}

func Unpack(enReq *model.EData) (*crypt.PData, error) {
	//privateKeyBytes, err := crypt.GetKeyBySdkVersion(enReq.SDKVersion)
	//if err != nil || privateKeyBytes == nil {
	//	logging.Info(err)
	//	return nil, err
	//}
	pData, err := crypt.Unpackv2(&crypt.EData{
		EK: []byte(enReq.EK),
		ED: []byte(enReq.ED),
		IV: []byte(enReq.IV),
	}, privateKeyBytes)
	if err != nil {
		logging.Info(err)
		return nil, err
	}
	return pData, nil
}

func UnClientPack(enReq *model.EDataResponse) ([]byte, error) {
	pData, err := crypt.UnClientPack([]byte(enReq.EK), []byte(enReq.ED), []byte(enReq.IV), privateKeyBytes)
	if err != nil {
		logging.Info(err)
		return nil, err
	}
	logging.Info(fmt.Sprintf("%s", pData))
	return pData, nil
}

//func bind(data string, v interface{}) error {
//	err := ffjson.Unmarshal([]byte(data), v)
//	if err != nil {
//		return err
//	}
//
//	err = validator.Validate.Struct(v)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
