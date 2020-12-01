package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"qlong/log"
	"strings"
	"thirdparty/tpwechatapi/cache"
	"thirdparty/tpwechatapi/utils"
	"time"
	"tools/action/constant/tpwechatapi"
	tutils "tools/utils"

	"gopkg.in/kataras/iris.v6"
)

// EventAppid 接收消息与事件
func EventAppid(ctx *iris.Context) {
	log.Info(">>>>>>>>>>>>>>>>>>>>>>>>> EventAppid 接收消息与事件")
	defer func() {
		if p := recover(); p != nil {
			log.Error(fmt.Sprintf("EventAppid | fun(%s) err(%+v)", ctx.Param("fun"), p))
			return
		}
	}()

	err := ctx.Request.ParseForm()
	if err != nil {
		utils.ReturnToClientString(ctx, "success")
		return
	}

	appid := ctx.Param("appid")
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> appid:%s\t APPID：%s", appid, tpwechatapi.APPID))
	// if appid == tpwechatapi.APPID {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.ReturnToClientString(ctx, "success")
		return
	}
	var t tpwechatapi.T
	t.Signature = strings.Join(ctx.Request.Form["signature"], "")
	t.Timestamp = strings.Join(ctx.Request.Form["timestamp"], "")
	t.Nonce = strings.Join(ctx.Request.Form["nonce"], "")
	t.MsgSignature = strings.Join(ctx.Request.Form["msg_signature"], "")
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> T：%+v", t))
	if t.MsgSignature != "" {
		var xt tpwechatapi.XT
		err = xml.Unmarshal(body, &xt)
		log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> XT：%+v", xt))
		// 构造对象
		wechatCryptor, _ := utils.NewWechatCryptor(tpwechatapi.APPID, tpwechatapi.TOKEN, tpwechatapi.ENCODINGAESKEY)
		xml, err := wechatCryptor.DecryptMsgContent(t.Signature, t.Timestamp, t.Nonce, xt.Encrypt)
		if CheckSignature(tpwechatapi.TOKEN, t.Signature, t.Timestamp, t.Nonce, t.MsgSignature) && err == nil {
			fmt.Printf("正确加解密\n")
			log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> auth/wechat/event 正确加解密：\n%s", xml))
			// 解密
			// xml, _ := wechatCryptor.DecryptMsg(t.Signature, t.Timestamp, t.Nonce, postData)
			// fmt.Printf("解密内容：%s\n", xml)
			// log.Error(fmt.Sprintf("解密内容：%s\n", xml))
			acceptMessageAndEvent(ctx, xml)
		}
	}
	// }
}

func acceptMessageAndEvent(ctx *iris.Context, xmll string) {
	var cb tpwechatapi.Callback
	err := xml.Unmarshal([]byte(xmll), &cb)
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> cb：%+v", cb))
	if err == nil {
		process(ctx, cb)
	}
}

func process(ctx *iris.Context, cb tpwechatapi.Callback) {
	if cb.MsgType == "event" {
		log.Info(">>>>>>>>>>>>>>>>>>>>>>>>> event")
		content := cb.Event + "from_callback"
		replyEventMessage(ctx, content, cb.ToUserName, cb.FromUserName)
	}
	if cb.MsgType == "text" {
		log.Info(">>>>>>>>>>>>>>>>>>>>>>>>> text")
		processTextMessage(ctx, cb.Content, cb.ToUserName, cb.FromUserName)
	}
}

func replyEventMessage(ctx *iris.Context, content, toUserName, fromUserName string) {
	createtime := fmt.Sprint(time.Now().Unix())
	xmlcontent := fmt.Sprintf(tpwechatapi.XmlContentFormat, fromUserName, toUserName, createtime, content)
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> xmlcontent：%s", xmlcontent))
	// 构造对象
	wechatCryptor, _ := utils.NewWechatCryptor(tpwechatapi.APPID, tpwechatapi.TOKEN, tpwechatapi.ENCODINGAESKEY)
	returnvalue, _ := wechatCryptor.EncryptMsg(xmlcontent, createtime, "easemob")
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> returnvalue：%s", returnvalue))
	utils.ReturnToClientString(ctx, xmlcontent)
	// utils.ReturnToClientString(ctx, returnvalue)
	return
}

func processTextMessage(ctx *iris.Context, content, toUserName, fromUserName string) {
	if content == "TESTCOMPONENT_MSG_TYPE_TEXT" {
		log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> TESTCOMPONENT_MSG_TYPE_TEXT"))
		returnContent := content + "_callback"
		replyEventMessage(ctx, returnContent, toUserName, fromUserName)
	} else if strings.HasPrefix(content, "QUERY_AUTH_CODE") || strings.HasPrefix(content, "query_auth_code") {
		log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> QUERY_AUTH_CODE"))
		utils.ReturnToClientString(ctx, "")
		if cache.Tools.GetCAT() == "" {
			resetAccessToken()
		}
		replyApiTextMessage(cache.Tools.GetCAT(), strings.Split(content, ":")[1], fromUserName)
	}
}

func resetAccessToken() {
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> resetAccessToken"))
	mp := make(map[string]string, 3)
	mp["component_appid"] = tpwechatapi.APPID
	mp["component_appsecret"] = tpwechatapi.APPSECRET
	mp["component_verify_ticket"] = cache.Tools.GetCVT()
	// 请求接口
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> GET_API_COMPONENT_TOKEN_URL：%s", tpwechatapi.GET_API_COMPONENT_TOKEN_URL))
	// str, err := tutils.HttpRequestPost(tpwechatapi.GET_API_COMPONENT_TOKEN_URL, mp)
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> resetAccessToken 请求参数：%+v", mp))
	str, err := tutils.HttpRequestPostJson(tpwechatapi.GET_API_COMPONENT_TOKEN_URL, mp)
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> str：%s", str))
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> err：%+v", err))
	var cat tpwechatapi.CAT
	err = json.Unmarshal([]byte(str), &cat)
	// log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> cat：%+v", cat))
	if err == nil && cat.Cat != "" {
		cache.Tools.SetCAT(cat.Cat, cat.Expires)
	}
}

func replyApiTextMessage(componentAccessToken, authCode, fromUserName string) {
	mp := make(map[string]string, 2)
	mp["component_appid"] = tpwechatapi.APPID
	mp["authorization_code"] = authCode
	// 请求接口
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> GET_API_QUERY_AUTH_URL：%s", tpwechatapi.GET_API_QUERY_AUTH_URL))
	str, err := tutils.HttpRequestPostJson(fmt.Sprintf(tpwechatapi.GET_API_QUERY_AUTH_URL, cache.Tools.GetCAT()), mp)
	// str, err := tutils.HttpRequestPost(fmt.Sprintf(tpwechatapi.GET_API_QUERY_AUTH_URL, cache.Tools.GetCAT()), mp)
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> str：%s", str))
	log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> err：%+v", err))
	type AuthorizationInfo struct {
		AuthorizationInfo struct {
			AuthorizerAppid        string `json:"authorizer_appid"`
			AuthorizerAccessToken  string `json:"authorizer_access_token"`
			ExpiresIn              int    `json:"expires_in"`
			AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
			FuncInfo               []struct {
				FuncscopeCategory struct {
					ID int `json:"id"`
				} `json:"funcscope_category"`
			} `json:"func_info"`
		} `json:"authorization_info"`
	}
	var ai AuthorizationInfo
	err = json.Unmarshal([]byte(str), &ai)
	// log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> ai：%+v", ai))
	if err == nil && ai.AuthorizationInfo.AuthorizerAccessToken != "" {
		mp := make(map[string]interface{}, 3)
		mpp := make(map[string]string, 1)
		mpp["content"] = authCode + "_from_api"
		mp["touser"] = fromUserName
		mp["msgtype"] = "text"
		mp["text"] = mpp
		log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> replyApiTextMessage 请求参数：%+v", mp))
		str, err := tutils.HttpRequestPostJson(fmt.Sprintf(tpwechatapi.SEND_URL, ai.AuthorizationInfo.AuthorizerAccessToken), mp)
		// str, err := tutils.HttpRequestPostJson(fmt.Sprintf(tpwechatapi.SEND_URL, cache.Tools.GetCAT()), mp)
		log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> str：%s", str))
		log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> err：%+v", err))
	}
}
