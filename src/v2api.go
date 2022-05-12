package main

import (
	"github.com/tidwall/gjson"
)

const V2API string = "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo"

func GetV2Quality(realRoomID string) JParam {
	param := JParam{"platform": "h5", "protocol": "1", "format": "0,1", "codec": "0", "room_id": realRoomID}
	return GetChooseQuality(param, "data.playurl_info.playurl.g_qn_desc", V2API)
}

func V2HandlerQualityUrl(quality int64, realRoomID string) JParam {
	param := JParam{"platform": "h5", "protocol": "1", "format": "0,1", "codec": "0", "room_id": realRoomID, "qn": quality}
	result := GetRequest(V2API, param)

	temp := gjson.Get(result, "data.playurl_info.playurl.stream.0.format.0.codec.0").String()
	baseUrl := gjson.Get(temp, "base_url").String()
	host := gjson.Get(temp, "url_info.0.host").String()
	extra := gjson.Get(temp, "url_info.0.extra").String()

	realUrl := [1]string{host + baseUrl + extra}

	return JParam{
		"urls": realUrl,
	}

}
