package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"goio/logger"
	"io/ioutil"
	"net/http"
)

type User struct {
	UserId int  `json:"userId"`
	Role   int  `json:"role"`
	Rm     bool `json:"remove"`
}

type Res struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data User   `json:"data"`
}

var (
	body []byte
	err  error
	res  Res
	resp *http.Response
)

func Post(url, ctype string, body []byte) (Res, error) {
	logger.Info("url %s Req body %s", url, string(body))
	resp, err = http.Post(url, ctype, bytes.NewBuffer(body))
	if err != nil {
		logger.Error("hp.Post error %v", err)
		return res, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	logger.Info("url %s Resp  body %s", url, string(body))
	if err != nil {
		logger.Error("hp.Post ioutil.ReadAll error %v", err)
		return res, err
	}

	err = json.Unmarshal(body, &res)

	if err != nil {
		logger.Error("hp.Post json.Unmarshal error %v", err)
		return res, err
	}
	return res, nil
}

func Get(url string) (Res, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("hp.Get error %v", err)
		return res, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("hp.Post ioutil.ReadAll error %v", err)
		return res, err
	}

	err = json.Unmarshal(body, &res)

	if err != nil {
		logger.Error("hp.Post json.Unmarshal error %v", err)
		return res, err
	}
	return res, err
}

func EncodeJson(notify interface{}) ([]byte, error) {
	return json.Marshal(notify)
}

func StoreMessage(rid string, body []byte) error {
	res, err := Post("http://"+GetHttpConfig().Remoteaddr+"/im/"+rid+"/chat", "application/json", body)
	if err != nil || res.Code != 0 {
		logger.Error("message store error %v", err)
		return errors.New("Store Message failed!")
	}
	return nil
}