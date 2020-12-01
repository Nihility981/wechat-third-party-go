package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"qlong/log"
	"strings"
	"thirdparty/tpwechatapi/cache"
	"thirdparty/tpwechatapi/utils"

	"gopkg.in/kataras/iris.v6"
)

// AuthEvent 接收授权事件
func AuthEvent(ctx *iris.Context) {
	defer func() {
		if p := recover(); p != nil {
			// log.Error(fmt.Sprintf("AuthEvent | fun(%s) err(%+v)", ctx.Param("fun"), p))
			utils.ReturnToClientString(ctx, "success")
			return
		}
	}()

	err := ctx.Request.ParseForm()
	if err != nil {
		utils.ReturnToClientString(ctx, "success")
		return
	}
	// log.Info(fmt.Sprintf("form:%v\n", ctx.Request.Form))
	// log.Info(fmt.Sprintf("body:%v\n", ctx.Request.Body))
	fmt.Printf("form:%v\n", ctx.Request.Form)
	fmt.Printf("body:%v\n", ctx.Request.Body)
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		// log.Error(fmt.Sprintf("ioutil.ReadAll(ctx.Request.Body) err：[ %+v ]", err))
		utils.ReturnToClientString(ctx, "success")
		return
	}
	var xt XT
	err = xml.Unmarshal(body, &xt)
	if err != nil {
		// log.Error(fmt.Sprintf("ioutil.ReadAll(ctx.Request.Body) err：[ %+v ]", err))
		utils.ReturnToClientString(ctx, "success")
		return
	}
	fmt.Printf("xt:%+v\n", xt)
	fmt.Printf("加密后: %s\n", xt.Encrypt)
	// log.Info(fmt.Sprintf("body:%+v\n", ctx.Request.Body))
	var t T
	t.Signature = strings.Join(ctx.Request.Form["signature"], "")
	t.Timestamp = strings.Join(ctx.Request.Form["timestamp"], "")
	t.Nonce = strings.Join(ctx.Request.Form["nonce"], "")
	t.MsgSignature = strings.Join(ctx.Request.Form["msg_signature"], "")
	// log.Info(fmt.Sprintf("t:%+v\n", t))
	utils.ReturnToClientString(ctx, "success")
	if t.MsgSignature != "" {
		// 处理
		processAuthorizeEvent(t, xt)
	}
}

func processAuthorizeEvent(t T, xt XT) {
	if t.MsgSignature == "" {
		return // 微信推送给第三方开放平台的消息一定是加过密的，无消息加密无法解密消息
	}
	// 构造对象
	wechatCryptor, _ := utils.NewWechatCryptor(APPID, TOKEN, ENCODINGAESKEY)
	xml, err := wechatCryptor.DecryptMsgContent(t.Signature, t.Timestamp, t.Nonce, xt.Encrypt)
	fmt.Printf("解密后明文: %s\n", xml)
	// log.Info(fmt.Sprintf("解密后明文: %s\n", xml))
	if CheckSignature(TOKEN, t.Signature, t.Timestamp, t.Nonce, t.MsgSignature) && err == nil {
		fmt.Printf("正确加解密\n")
		log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>>>>>> auth/event 正确加解密：\n%s", xml))
		// 解密
		// xml, _ := wechatCryptor.DecryptMsg(t.Signature, t.Timestamp, t.Nonce, postData)
		processAuthorizationEvent(xml)
	}
}

func processAuthorizationEvent(xmll string) {
	var cvt CVT
	err := xml.Unmarshal([]byte(xmll), &cvt)
	if err != nil {
		log.Error(fmt.Sprintf("processAuthorizationEvent err：[ %v ]", err))
		return
	}
	fmt.Printf("cvt:%+v\n", cvt)
	fmt.Printf("ComponentVerifyTicket: %s\n", cvt.ComponentVerifyTicket)
	if cvt.InfoType == "component_verify_ticket" && cvt.ComponentVerifyTicket != "" {
		cache.Tools.SetCVT(cvt.ComponentVerifyTicket)
	}
}

func CheckSignature(tk, signature, timestamp, nonce, encrypt string) bool {
	fmt.Printf("检查是否正确加密：[tk(%s) signature(%s) timestamp(%s) nonce(%s)]\n", tk, signature, timestamp, nonce)
	if tk != "" && signature != "" && timestamp != "" && nonce != "" {
		sha1 := utils.SHA1(tk, timestamp, nonce, encrypt)
		fmt.Printf("sha1：%s\n", sha1)
		if sha1 != signature {
			fmt.Printf("错误加密\n")
			return false
		}
		return true
	}
	return false
}
