package main

import (
	"fmt"
	"qlong"
	"thirdparty/tpwechatapi/cache"
	"tools/action/constant/publicconfig"
)

const (
	APPID                       = ""
	TOKEN                       = ""
	ENCODINGAESKEY              = ""
	APPSECRET                   = ""
	GET_API_COMPONENT_TOKEN_URL = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	GET_API_QUERY_AUTH_URL      = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=%s"
	SEND_URL                    = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s"
)

type T struct {
	Nonce        string `json:"nonce"`
	Timestamp    string `json:"timestamp"`
	Signature    string `json:"signature"`
	MsgSignature string `json:"msg_signature"`
}

type XT struct {
	Encrypt      string `xml:"Encrypt"`
	MsgSignature string `xml:"MsgSignature"`
	TimeStamp    string `xml:"TimeStamp"`
	Nonce        string `xml:"Nonce"`
}

type CVT struct {
	AppId                 string `xml:"AppId"`
	CreateTime            string `xml:"CreateTime"`
	InfoType              string `xml:"InfoType"`
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
}

type Callback struct {
	MsgType      string `xml:"MsgType"`
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	Event        string `xml:"Event"`
	Content      string `xml:"Content"`
}

type CAT struct {
	Cat     string `json:"component_access_token"`
	Expires int64  `json:"expires_in"`
}

const XmlContentFormat = `
<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%s</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>
`

func main() {
	// 初始化qlong
	qlong.NewQLong("config.conf")
	fmt.Println("端口号:", qlong.Port)
	if qlong.Atag == "debug" {
		qlong.Atag = "tpwechatapi"
		qlong.Port = "6006"
	}
	// 连接redis
	cacheErr := cache.Tools.Init(publicconfig.RedisHostBiz, publicconfig.RedisLatePush)
	if cacheErr != nil {
		qlong.LogError("main", "redis", "", fmt.Sprintf("初始化redis失败:%v", cacheErr))
		return
	}

	//启动iris,并监听
	app := qlong.GetIris(-1)
	app.Any("/auth/event", AuthEvent)
	app.Any("/auth/wechat/event/:appid", EventAppid)
	app.Listen(":" + qlong.Port)

}
