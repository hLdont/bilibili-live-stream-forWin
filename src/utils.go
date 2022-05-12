package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
)

type JParam map[string]interface{}

func (this JParam) set(key string, value interface{}) {
	this[key] = value
}

func (this JParam) toUrlQuery() string {
	paramsTemp := url.Values{}
	for k, v := range this {
		paramsTemp.Set(k, fmt.Sprintf("%v", v))
	}
	return paramsTemp.Encode()
}

func (this JParam) toJson() string {
	marshal, _ := json.Marshal(this)
	return string(marshal)
}

func (this JParam) Println() {
	marshal, _ := json.Marshal(this)
	fmt.Println(string(marshal))
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}

func GetChooseQuality(param JParam, path string, api string) JParam {
	result := GetRequest(api, param)

	var qualities []JParam
	json.Unmarshal([]byte(gjson.Get(result, path).Raw), &qualities)
	return JParam{
		"quality": qualities,
	}
}

func WriteString(content string) {
	fileName := "urls.txt"
	var dstFile *os.File
	if !IsExists(fileName) {
		dstFile, _ = os.Create(fileName)
	} else {
		_ = os.Remove(fileName)
		dstFile, _ = os.Create(fileName) // easy way to io use
	}

	defer func(dstFile *os.File) {
		_ = dstFile.Close()
	}(dstFile)

	_, _ = dstFile.WriteString(content)
}

func GetRequest(address string, params JParam) string {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(JParam{
				"type": -1,
				"data": err,
			})
		}
	}()

	Url, _ := url.Parse(address)
	Url.RawQuery = params.toUrlQuery()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", Url.String(), nil)

	res, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)

	return string(body)

}
