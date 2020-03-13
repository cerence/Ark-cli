package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var UnknowError = errors.New("request fail, please contact arkcloud@cerence.com")

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

func ImportNewDevice(host string, uniqueDeviceCode string) error {
	deviceCode, vin, err := splitUniqueCode(uniqueDeviceCode)
	if err != nil {
		return err
	}

	req := make([]ImportReq, 0)
	req = append(req, ImportReq{
		DeviceCode:  deviceCode,
		VIN:         vin,
		Make:        "Ark Edge Stub",
		Model:       "Ark Edge Stub",
		ModelYear:   "2020",
		TrimLevel:   "SE",
		FuelType:    "Petrol",
		HasHeadUnit: "true",
		NMAID:       "test",
	})

	url := fmt.Sprintf("http://%s/api/devices/import", host)
	var body []byte
	body, err = json.Marshal(req)
	if err != nil {
		return err
	}
	if err := HttpRequestUtil(url, "POST", body, nil); err != nil {
		return UnknowError
	} else {
		Info.Println("Import complete")
	}

	return nil
}
func UnbindDevice(host string, uniqueDeviceCode string) error {
	device, err := queryDevice(host, uniqueDeviceCode)
	if err != nil {
		return err
	}

	if device.ID == 0 {
		return errors.New("the device id is wrong")
	}

	Info.Println(fmt.Sprintf("The device id is %d", device.ID))

	binding := &DeviceBinding{
		DeviceID:  device.ID,
		IsBinding: false,
	}
	var body []byte
	body, err = json.Marshal(binding)
	if err != nil {
		return err
	}
	res := &Response{}
	url := fmt.Sprintf("http://%s/api/devices/binding", host)
	if err := HttpRequestUtil(url, "PUT", body, res); err != nil {
		return UnknowError
	}
	if res.ErrCode > 0 {
		return errors.New(fmt.Sprintf("query device error: %s", res.ErrMsg))
	}

	Info.Println("Unbind device complete")
	return nil
}

func queryDevice(host string, uniqueDeviceCode string) (*Device, error) {
	url := fmt.Sprintf("http://%s/api/devices/info?unique_device_code=%s", host, uniqueDeviceCode)
	res := &Response{}

	if err := HttpRequestUtil(url, "GET", nil, res); err != nil {
		return nil, err
	}

	if res.ErrCode > 0 {
		return nil, errors.New(fmt.Sprintf("query device error: %s", res.ErrMsg))
	}

	data, err := json.Marshal(res.ResData)
	if err != nil {
		return nil, err
	}
	devices := make([]Device, 0)
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, errors.New("no such device")
	}

	return &devices[0], nil
}

func splitUniqueCode(uniqueDeviceCode string) (string, string, error) {
	temp := strings.Split(uniqueDeviceCode, "_")
	if len(temp) != 2 || temp[0] == "" || temp[1] == "" {
		return "", "", errors.New("invalid uniqueDeviceCode")
	}

	return temp[0], temp[1], nil
}
