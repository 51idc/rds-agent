package funcs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/51idc/rds-agent/g"
)

type smartAPI_Data struct {
	Endpoint string `json:"endpoint"`
	Version  string `json:"version"`
}

type smartAPI_Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func sendData(url string, data smartAPI_Data) ([]byte, int, error) {

	js, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}
	res, err := http.Post(url, "Content-Type: application/json", bytes.NewBuffer(js))
	if err != nil {
		return nil, 0, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}
	defer res.Body.Close()
	return body, res.StatusCode, err
}

func smartAPI_Push(version string) {
	debug := g.Config().Debug
	smartAPI_url := g.Config().SmartAPI.Url
	endpoint, err := g.Hostname()
	var data smartAPI_Data
	var result smartAPI_Result

	data.Endpoint = endpoint
	data.Version = version
	body, res, err := sendData(smartAPI_url, data)
	if err != nil {
		log.Println(err)
		return
	}
	if res != 200 {
		log.Println("smartAPI error,statcode= ", res)
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return
	}
	if result.Status != "ok" {
		log.Println("SmartAPI return error: ", result.Message)
		return
	}
	if debug {
		log.Println("Push Version to SmartAPI Success: ", version)
	}
	return
}
