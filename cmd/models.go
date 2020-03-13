package cmd

import "time"

type Response struct {
	ErrCode int         `json:"error_code"`
	ErrMsg  string      `json:"error_msg"`
	ResData interface{} `json:"response_data"`
}

type ImportReq struct {
	DeviceCode  string `json:"DeviceCode"`
	VIN         string `json:"VIN"`
	Make        string `json:"Make"`
	Model       string `json:"Model"`
	ModelYear   string `json:"ModelYear"`
	TrimLevel   string `json:"TrimLevel"`
	FuelType    string `json:"FuelType"`
	HasHeadUnit string `json:"HasHeadUnit"`
	NMAID       string `json:"NMAID"`
}

type Device struct {
	ID               int       `json:"id"`
	UniqueDeviceCode string    `json:"unique_device_code"`
	NickName         string    `json:"nick_name"`
	DeviceModelID    int       `json:"device_model_id"`
	Description      string    `json:"description"`
	FwVersion        string    `json:"fw_version"`
	IsRetired        string    `json:"is_retired"`
	DeviceCode       string    `json:"device_code"`
	AttachDeviceID   string    `json:"attach_device_id"`
	HasHU            string    `json:"has_hu"`
	NMAID            string    `json:"nmaid"`
	ImportDate       time.Time `json:"import_date"`
}

type DeviceBinding struct {
	DeviceID  int    `json:"device_id"`
	IsBinding bool   `json:"is_binding"`
	UserName  string `json:"user_name"`
}

type DeviceRegReq struct {
	DeviceCode string `json:"device_code" binding:"required"`
	VIN        string `json:"VIN" binding:"required"`
	NMAID      string `json:"nmaid" binding:"required"`
}

type TokenResData struct {
	Token    string `json:"token"`
	ExpireIn int    `json:"expire_in"`
}
