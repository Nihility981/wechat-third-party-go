package utils

//ReturnAck 返回错误代码
import (
	"fmt"
	"qlong"
	"qlong/utl"
	"sort"
	"strconv"
	"strings"
	"time"

	mathrand "math/rand"

	iris "gopkg.in/kataras/iris.v6"
)

//ReturnYunNanAck 返回错误代码
func ReturnYunNanAck(ctx *iris.Context, ack, msg string) {
	var m1 map[string]interface{}
	m1 = make(map[string]interface{})
	intack := ack
	m1["error_code"] = intack
	m1["error_msg"] = msg
	ReturnToClient(ctx, m1)
}

//ReturnHuNanAck 返回错误代码
func ReturnHuNanAck(ctx *iris.Context, ack, msg string) {
	var m1 map[string]interface{}
	m1 = make(map[string]interface{})
	intack := ack
	m1["status"] = intack
	m1["error_msg"] = msg
	ReturnToClient(ctx, m1)
}

// ReturnToClient 设置跨域访问
func ReturnToClient(ctx *iris.Context, m1 map[string]interface{}) {
	// ctx.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Header().Set("Access-Control-Allow-Origin", ctx.RequestHeader("Origin"))
	ctx.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type")
	ctx.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx.JSON(iris.StatusOK, m1)
}

// 设置跨域访问
func ReturnToClientString(ctx *iris.Context, rtn string) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type")
	ctx.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Header().Set("Content-Type", "text/plain")
	ctx.WriteString(rtn)
}

//根据参数map获取签名
func GetAuthSign(params map[string]string, skey string) string {
	sn := ""
	//取出所有参数名
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	//对参数名做升序排列
	sort.Strings(keys)
	strs := ""
	for inx, key := range keys {
		val := params[key]
		if val == "" || val == "null" {
			continue
		}
		if inx == 0 {
			strs += key + "=" + val
		} else {
			strs += "&" + key + "=" + val
		}
	}
	strs += "&SecretKey=" + skey
	qlong.DBG("拼装后的签名字符串为[%s]", strs)
	sn = utl.MD5(strs)
	sn = strings.ToLower(sn)
	qlong.DBG("获取到的签名为[%s]", sn)
	return sn
}

//生成请求头签名
func GetSign(params map[string]string) string {
	sn := ""
	//取出所有参数名
	var keys []string
	for _, val := range params {
		keys = append(keys, val)
	}
	//对参数值做升序排列
	sort.Strings(keys)
	strs := ""
	for _, key := range keys {
		strs += key
	}
	qlong.DBG("2拼装后的签名字符串为[%s]", strs)
	sn = utl.MD5(strs)
	sn = strings.ToLower(sn)
	qlong.DBG("2获取到的签名为[%s]", sn)
	return sn
}

//生成请求头验证信息
func GetAuthorization(appid, tm, nonceStr string) string {
	appid = strings.ToLower(appid)
	strs := appid + ":" + tm + ":" + nonceStr
	return utl.Base64Encode(strs)
}

//验证信息反解析出AppID + 冒号 + timestamp + 冒号 + nonceStr
func DecodeAuthorization(str string) (string, string, string, bool) {
	str = utl.Base64Decode(str)
	strs := strings.Split(str, ":")
	if len(strs) == 3 {
		return strs[0], strs[1], strs[2], true
	}
	return "", "", "", false
}

//生成长度24位的唯一数字单号
func MakeYearDaysRand() string {
	strs := time.Now().Format("06")
	days := strconv.Itoa(GetDaysInYearByThisYear())
	count := len(days)
	if count < 3 {
		days = strings.Repeat("0", 3-count) + days
	}
	strs += days
	sum := 19
	var untime = time.Now().UnixNano()
	var keys = mathrand.Intn(int(untime)) + int(untime)
	mathrand.Seed(int64(keys))
	result := strconv.Itoa(mathrand.Intn(int(untime)))
	count = len(result)
	if count < sum {
		result = strings.Repeat("0", sum-count) + result
	}
	strs += result
	if len(strs) > 24 {
		strs = string([]rune(strs)[:24])
	}
	return strs
}

func GetDaysInYearByThisYear() int {
	now := time.Now()
	total := 0
	arr := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	y, month, d := now.Date()
	m := int(month)
	for i := 0; i < m-1; i++ {
		total = total + arr[i]
	}
	if (y%400 == 0 || (y%4 == 0 && y%100 != 0)) && m > 2 {
		total = total + d + 1

	} else {
		total = total + d
	}
	return total
}

// 将yyyy-MM-dd HH:mm:ss 格式化成 20190918  175642 返回
func AFormatDateTime(datestr string) (string, string) {

	if datestr == "" || len(datestr) < 19 {
		return "", ""
	}
	date := datestr[0:4] + datestr[5:7] + datestr[8:10]
	time := datestr[11:13] + datestr[14:16] + datestr[17:19]

	return date, time
}

// GetCurrentTime 得到当前时间
func GetCurrentTime(iFlag int) string {
	if iFlag == 0 { //获得当前日期 yyyy-MM-dd HH:mm:ss
		return time.Now().Format("2006-01-02 15:04:05")
	} else if iFlag == 1 { //获得当前日期 yyyyMMdd
		return time.Now().Format("20060102")
	} else if iFlag == 2 { //获得当前日期 HHmmss
		return time.Now().Format("150405")
	} else if iFlag == 3 { //获得当前年月 HHmm
		return time.Now().Format("200601")
	} else if iFlag == 4 { //获得当前年月日时分秒
		return time.Now().Format("20060102150405")
	} else if iFlag == 5 { //获得当前年
		return time.Now().Format("2006")
	} else if iFlag == 6 { //获取当前月日
		return time.Now().Format("0102")
	} else if iFlag == 7 { //获取时间戳,精确到纳秒
		return fmt.Sprintf("%d", time.Now().UnixNano())
	} else if iFlag == 8 { //获取时间戳,精确到毫秒
		return fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	}
	return time.Now().Format("20060102150405")
}
