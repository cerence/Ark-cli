package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var UnknownError = errors.New("request fail, please contact arkcloud@cerence.com")
var NoDeviceError = errors.New("no such device")
var DeviceCodeSplitter = "_"

func ImportNewDevice(apiHost string, backendHost string, uniqueDeviceCode string) error {
	deviceCode, vin, err := splitUniqueCode(uniqueDeviceCode)
	if err != nil {
		return err
	}

	device, err := queryDevice(backendHost, uniqueDeviceCode)
	if err != NoDeviceError || (device != nil && device.ID > 0) {
		return errors.New("the device is existed")
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

	url := fmt.Sprintf("http://%s/api/devices/import", backendHost)
	var body []byte
	body, err = json.Marshal(req)
	if err != nil {
		return err
	}
	if err := HttpRequestUtil(url, "POST", body, nil); err != nil {
		return UnknownError
	}

	token, err := registerDevice(apiHost, deviceCode, vin)
	if err != nil {
		return err
	}

	Info.Println("Import complete")
	Info.Println(fmt.Sprintf("Token is %s", token))
	return nil
}

func UnbindDevice(backendHost string, uniqueDeviceCode string) error {
	device, err := queryDevice(backendHost, uniqueDeviceCode)
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
	url := fmt.Sprintf("http://%s/api/devices/binding", backendHost)
	if err := HttpRequestUtil(url, "PUT", body, res); err != nil {
		return UnknownError
	}
	if res.ErrCode > 0 {
		return errors.New(fmt.Sprintf("query device error: %s", res.ErrMsg))
	}

	Info.Println("Unbind device complete")
	return nil
}

func registerDevice(apiHost string, deviceCode string, vin string) (string, error) {
	url := fmt.Sprintf("http://%s/v1/devices/%s/registration", apiHost, deviceCode+DeviceCodeSplitter+vin)
	req := &DeviceRegReq{
		DeviceCode: deviceCode,
		VIN:        vin,
		NMAID:      "test",
	}
	var body []byte
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	res := &Response{}
	if err := HttpRequestUtil(url, "POST", body, res); err != nil {
		return "", UnknownError
	}

	if res.ErrCode > 0 {
		return "", errors.New(res.ErrMsg)
	}
	tokenRes := &TokenResData{}
	if err := getResData(res.ResData, tokenRes); err != nil {
		return "", err
	}

	return tokenRes.Token, nil
}

func queryDevice(backendHost string, uniqueDeviceCode string) (*Device, error) {
	url := fmt.Sprintf("http://%s/api/devices/info?unique_device_code=%s", backendHost, uniqueDeviceCode)
	res := &Response{}

	if err := HttpRequestUtil(url, "GET", nil, res); err != nil {
		return nil, err
	}

	if res.ErrCode > 0 {
		return nil, errors.New(fmt.Sprintf("query device error: %s", res.ErrMsg))
	}

	devices := make([]Device, 0)
	if err := getResData(res.ResData, &devices); err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, NoDeviceError
	}

	return &devices[0], nil
}

func getResData(raw interface{}, result interface{}) error {
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	return nil
}

func splitUniqueCode(uniqueDeviceCode string) (string, string, error) {
	temp := strings.Split(uniqueDeviceCode, DeviceCodeSplitter)
	if len(temp) != 2 || temp[0] == "" || temp[1] == "" {
		return "", "", errors.New("invalid uniqueDeviceCode")
	}

	return temp[0], temp[1], nil
}
