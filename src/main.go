package main

import (
	"flag"
	"strconv"

	"github.com/tidwall/gjson"
)

const URL string = "https://api.live.bilibili.com/xlive/web-room/v1/playUrl/playUrl"

var apiSelect int
var typeExec int
var roomID string
var qualitySelect int64

func initArgs() {
	flag.IntVar(&apiSelect, "apiType", 1, "1: v1 api\n 2: v2 api \n")
	flag.IntVar(&typeExec, "type", 0, "0: get room real id and stream quality\n 1: get real url\n")
	flag.StringVar(&roomID, "id", "-1", "get room id with param")
	flag.Int64Var(&qualitySelect, "quality", 0, "select quality")
	flag.Parse()

	if apiSelect < 1 || apiSelect > 2 {
		panic("接口选择apiType 非法")
	} else if roomId_int, err := strconv.Atoi(roomID); err != nil || roomId_int < 0 {
		panic("房间号 非法")
	} else if typeExec < 0 || typeExec > 2 {
		panic("接口类型type 非法")
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			JParam{
				"type": -1,
				"data": err,
			}.Println()
		}
	}()
	initArgs()
	switch typeExec {
	case 0:
		realRoomID := GetRealRoomId()
		JParam{
			"type":    0,
			"realId":  realRoomID,
			"Quality": GetQuality(realRoomID),
		}.Println()
	case 1:
		GetRealUrl(roomID, qualitySelect).Println()
	default:
		return
	}

}

// case 0: get room real id
func GetRealRoomId() string {
	// 如果直播间未开播/轮播有bug，没改
	roomResult := GetRequest("https://api.live.bilibili.com/room/v1/Room/room_init", JParam{"id": roomID})
	RealId := handlerLiveStatus(roomResult)
	return RealId
}

func handlerLiveStatus(result string) string {
	code := gjson.Get(result, "code").Int()
	if code != 0 {
		if code == 60004 {
			panic("直播间不存在")
		}
		liveStatus := gjson.Get(result, "data.live_status").Int()
		if liveStatus != 1 {
			panic("直播间未开播")
		}
		panic("未知错误")
	} else {
		// 如果直播间未开播/轮播有bug，code照样是0，没改
		return gjson.Get(result, "data.room_id").String()
	}
}

// case 0: get stream quality with real room id
func GetQuality(realRoomID string) JParam {
	if apiSelect == 1 {
		return GetV1Quality(realRoomID)
	} else if apiSelect == 2 {
		return GetV2Quality(realRoomID)
	}
	panic("api 不存在")
}

// case 1: get stream url
func GetRealUrl(realRoomID string, qualitySelect int64) JParam {
	if apiSelect == 1 {
		return V1HandlerQualityUrl(qualitySelect, realRoomID)
	} else if apiSelect == 2 {
		return V2HandlerQualityUrl(qualitySelect, realRoomID)
	}
	panic("api 不存在")
}
