package model

type EData struct {
	SDKVersion string `json:"sv" binding:""`
	EK         string `json:"ek" binding:"required"`
	IV         string `json:"iv" binding:"required"`
	ED         string `json:"ed" binding:"required"`
}

type EDataResponse struct {
	EK string `json:"ek"`
	IV string `json:"iv"`
	ED string `json:"ed"`
}
