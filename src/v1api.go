package main

import (
	"github.com/tidwall/gjson"
)

const V1API string = "https://api.live.bilibili.com/xlive/web-room/v1/playUrl/playUrl"

func GetV1Quality(realRoomID string) JParam {
	param := JParam{"cid": realRoomID, "platform": "h5"}
	return GetChooseQuality(param, "data.quality_description", V1API)
}

func V1HandlerQualityUrl(qn int64, realRoomID string) JParam {
	param := JParam{"cid": realRoomID, "platform": "h5", "qn": qn}
	result := GetRequest(V1API, param)

	var urls []string

	gjson.Get(result, "data.durl").ForEach(func(key, value gjson.Result) bool {
		value.Get("url").ForEach(func(key, value gjson.Result) bool {
			urls = append(urls, value.String())
			return true
		})
		return true
	})

	return JParam{
		"urls": urls,
	}
}
