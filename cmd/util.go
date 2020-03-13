package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func HttpRequestUtil(url string, method string, request []byte, result interface{}) error {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return errors.New("bad request")
	}

	if result != nil {
		if err = json.Unmarshal(resBody, result); err != nil {
			return err
		}
	}

	return nil
}

var (
	Trace   *log.Logger // Just about anything
	Info    *log.Logger // Important information
	Warning *log.Logger // Be concerned
	Error   *log.Logger // Critical problem
)

func init() {

	Trace = log.New(os.Stdout, "TRACE: ", log.Lshortfile)

	Info = log.New(os.Stdout, "INFO: ", log.Lshortfile)

	Warning = log.New(os.Stdout, "WARNING: ", log.Lshortfile)

	Error = log.New(io.MultiWriter(os.Stdout, os.Stderr), "ERROR: ", log.Lshortfile)
}
