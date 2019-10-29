package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"jcd/service/dbcomm"
	"jcd/service/smscode"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mojocn/base64Captcha"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"bytes"
	"encoding/base64"
	"image/png"
	"io/ioutil"
	"os"
)

var (
	WX_PAY_URL   = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	WX_QUERY_URL = "https://api.mch.weixin.qq.com/pay/orderquery"

	CPCN_QUERY_URL = "http://192.168.70.180:8080/fcp2intra/light/life/query?data="
	CPCN_PAY_URL   = "http://192.168.70.180:8080/fcp2intra/light/life?data="

	WX_PAY_CALLBACK_URL = "http://www.doulaikan.club/jc/api/wxpaycallback"
	MCT_ID              = "1452819402"
	APP_ID              = "wx2db791be2eb77467"
	PRODUCT_NAME        = "测试商品"
	SERVER_IP           = "132.232.11.85"
	MCT_KEY             = "qj837vwk83xk2902jidk93slw82ms5ka"
	TRADE_TYPE_NATIVE   = "NATIVE"
	TRADE_TYPE_JSAPI    = "JSAPI"
	//=================
	WEIBO_OAUTH_CALLBACK_URL  = "http://www.doulaikan.club/jc/api/wxcallback"
	QQ_OAUTH_CALLBACK_URL     = "http://www.doulaikan.club/jc/api/wxcallback"
	WECHAT_OAUTH_CALLBACK_URL = "http://www.doulaikan.club/jc/api/wxcallback"

	CRF_SMSURL_UAT = "http://192.168.70.196:8080/sms_platform/receiveController/"
	CRF_SMSURL_PRD = "http://shortmessage-uat.crfchina.com/sms_platform/receiveController/"

	SMS_POST_KEY     = "f2dd8a0cf33c4f448b6dff8d13e98ae9"
	SMS_SYS_NO       = "36"
	SMS_SERVICE_TYPE = "crf_xiaocui_001"
)

const ERR_CODE_SUCCESS = "0000"
const ERR_CODE_DBERROR = "1001"
const ERR_CODE_TOKENER = "1003"
const ERR_CODE_PARTOEN = "1005"
const ERR_CODE_JSONERR = "2001"
const ERR_CODE_URLERR = "2005"
const ERR_CODE_NOTFIND = "3000"
const ERR_CODE_NOMATCH = "3010"
const ERR_CODE_EXPIRED = "6000"
const ERR_CODE_TYPEERR = "4000"
const ERR_CODE_STATUS = "5000"
const ERR_CODE_FAILED = "9000"
const ERR_CODE_OPERTYP = "4005"
const ERR_CODE_EXISTED = "4040"
const ERR_CODE_TOOBUSY = "6010"
const ERR_CODE_VERIFY = "7020"
const ERR_CODE_PAYERR = "8010"
const ERR_CODE_QRCODE = "7060"
const ERR_CODE_PAYDOING = "6666"

const ERR_USER_MSTSIGNUP = "7901"
const ERR_USER_SIGNINED = "7902"
const ERR_USER_UNSIGNIN = "7903"
const ERR_SMS_SENDERR = "9066"

const LOGIN_PHONE = 1
const LOGIN_OAUTH = 2

const STATUS_DISABLED = 1
const STATUS_ENABLED = 0
const STATUS_SUCC = "S"
const STATUS_INIT = "I"
const STATUS_FAIL = "F"
const STATUS_DOING = "D"

const MAX_SEARCH_TIMES = "5"

const FIELD_LOGIN_PASS = "login_pass"
const FIELD_ERRORS = "errors"
const FIELD_KILLS = "kills"
const FIELD_LIKES = "likes"
const FIELD_UPDATE_TIME = "update_time"
const FIELD_UPDATE_USER = "update_user"
const FIELD_PROC_STATUS = "proc_status"
const FIELD_PROC_MSG = "proc_msg"
const FIELD_PREPAY_ID = "prepay_id"
const FIELD_CODE_URL = "code_url"

const DEFAULT_PWD = "123456"
const SUCC_MSG = "success"
const EMPTY_STRING = ""

const SMSTYPE_LOGIN = "login"
const SMSTYPE_RESET = "reset"
const SMS_STATUS_INIT = "i"
const SMS_STATUS_END = "e"

const SMSCODE_EXPIRED_MINUTE = 20
const SMSCODE_MIN_INTERVAL = 10

const COMMENT_INIT_VALUE = 0
const COMMENT_LIKE = 10
const COMMENT_KILL = 20
const COMMENT_REPLY = 30

const REGION_PROVINCE = 1
const REGION_CITY = 2
const REGION_COUNTY = 3
const REGION_TOWN = 4
const REGION_VILLAGE = 5

var (
	ERROR_MAP map[string]string = map[string]string{
		ERR_CODE_SUCCESS:   "执行成功:",
		ERR_CODE_DBERROR:   "DB执行错误:",
		ERR_CODE_JSONERR:   "JSON格式错误:",
		ERR_CODE_EXPIRED:   "时效已经到期:",
		ERR_CODE_TYPEERR:   "类型转换错误:",
		ERR_CODE_STATUS:    "状态不正确:",
		ERR_CODE_TOKENER:   "获取TOKEN失败:",
		ERR_CODE_PARTOEN:   "解析TOKEN错误:",
		ERR_CODE_NOMATCH:   "比较不匹配:",
		ERR_CODE_URLERR:    "Url传参有误:",
		ERR_CODE_OPERTYP:   "ShowType类型错误:",
		ERR_CODE_NOTFIND:   "查询没发现提示:",
		ERR_CODE_EXISTED:   "注册账户已经存在:",
		ERR_CODE_TOOBUSY:   "短信发送太频繁:",
		ERR_CODE_VERIFY:    "验证码校验错误:",
		ERR_CODE_PAYERR:    "支付交易失败:",
		ERR_CODE_QRCODE:    "生产支付扫描失败:",
		ERR_USER_MSTSIGNUP: "用户没注册，需要注册",
		ERR_USER_SIGNINED:  "用户已经登录",
		ERR_USER_UNSIGNIN:  "用户需要登录",
		ERR_SMS_SENDERR:    "发送短信出错",
		ERR_CODE_PAYDOING:  "支付中",
	}
)

type ErrorResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

const (
	USER_CHARGE = "用户充值"
	FLOW_CHARGE = "charge"
	FLOW_INIT   = "i"
	FLOW_SUCC   = "s"
	FLOW_FAIL   = "f"

	NOW_TIME_FORMAT    = "2006-01-02 15:04:05"
	FIELD_ACCOUNT_BAL  = "Account_bal"
	FIELD_UPDATED_TIME = "Updated_time"

	CODE_SUCC    = "0000"
	CODE_NOEXIST = "1000"

	CODE_FAIL = "2000"

	RESP_SUCC = "0000"
	RESP_FAIL = "1000"

	CODE_TYPE_EDU       = "EDU"
	CODE_TYPE_POSITION  = "POSITION"
	CODE_TYPE_SALARY    = "SALARY"
	CODE_TYPE_WORKYEARS = "WORKYEARS"
	CODE_TYPE_POSICLASS = "POSICLASS"
	CODE_TYPE_REWARDS   = "REWARDS"

	TOKEN_KEY = "u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4"
)

func PrintHead(a ...interface{}) {
	log.Println("========》", a)
}

func PrintTail(a ...interface{}) {
	log.Println("《========", a)
}

func Write_Response(response interface{}, w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "1728000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,Action, Module,Authorization")
	fmt.Fprintf(w, string(json))
}

/*
	说明：用户登录成功后，产生SESSION的TOKEN
	入参：
	出参：参数1：Token
		 参数1：Error
*/

func GetToken(userId string, nickName string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["aud"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(10)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["cnt"] = MAX_SEARCH_TIMES
	claims["nne"] = nickName

	token.Claims = claims
	tokenString, err := token.SignedString([]byte(TOKEN_KEY))
	if err != nil {
		return EMPTY_STRING, err
	}
	return tokenString, nil
}

/*
	说明：根据TOKEN 进行校验
	入参：
	出参：参数1：Token
		 参数1：Error
*/

func CheckToken(w http.ResponseWriter, req *http.Request) (string, string, string, error) {
	PrintHead("CheckToken")
	var errResp ErrorResp
	auth := req.Header.Get("Authorization")
	var authB64 string
	auths := strings.SplitN(auth, " ", 2)
	if len(auths) != 2 {
		authB64 = auth
	} else {
		authB64 = auths[1]
	}
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(authB64, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(TOKEN_KEY), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			errResp.ErrCode = ERR_CODE_EXPIRED
		} else {
			errResp.ErrCode = ERR_CODE_PARTOEN
		}
		errResp.ErrMsg = ERROR_MAP[ERR_CODE_PARTOEN] + err.Error()
		Write_Response(errResp, w, req)
		return EMPTY_STRING, EMPTY_STRING, EMPTY_STRING, err
	}
	userId, ok := claims["aud"].(string)
	if !ok {
		errResp.ErrCode = ERR_CODE_TYPEERR
		errResp.ErrMsg = ERROR_MAP[ERR_CODE_TYPEERR] + "userid"
		Write_Response(errResp, w, req)
		return EMPTY_STRING, EMPTY_STRING, EMPTY_STRING, fmt.Errorf("Assertion Error.")
	}
	maxTimes, ok := claims["cnt"].(string)
	if !ok {
		errResp.ErrCode = ERR_CODE_TYPEERR
		errResp.ErrMsg = ERROR_MAP[ERR_CODE_TYPEERR] + "maxTimes"
		Write_Response(errResp, w, req)
		return EMPTY_STRING, EMPTY_STRING, EMPTY_STRING, fmt.Errorf("Assertion Error.")
	}
	nickName, ok := claims["nne"].(string)
	if !ok {
		errResp.ErrCode = ERR_CODE_TYPEERR
		errResp.ErrMsg = ERROR_MAP[ERR_CODE_TYPEERR] + "nickName"
		Write_Response(errResp, w, req)
		return EMPTY_STRING, EMPTY_STRING, EMPTY_STRING, fmt.Errorf("Assertion Error.")
	}

	return userId, maxTimes, nickName, nil
}

/*
	说明：根据TOKEN 进行校验
	入参：
	出参：参数1：Token
		 参数1：Error
*/

func CheckTokenExt(w http.ResponseWriter, req *http.Request) (string, string, string, error) {
	PrintHead("CheckTokenExt")
	auth := req.Header.Get("Authorization")
	var authB64 string
	auths := strings.SplitN(auth, " ", 2)
	if len(auths) != 2 {
		authB64 = auth
	} else {
		authB64 = auths[1]
	}
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(authB64, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(TOKEN_KEY), nil
	})
	if err != nil {
		//if err.Error() == "Token is expired" {
		return EMPTY_STRING, EMPTY_STRING, EMPTY_STRING, err
	}
	userId, ok := claims["aud"].(string)
	if !ok {

		return EMPTY_STRING, EMPTY_STRING, EMPTY_STRING, fmt.Errorf("Assertion Error.")
	}
	maxTimes, ok := claims["cnt"].(string)
	if !ok {
		return EMPTY_STRING, EMPTY_STRING, EMPTY_STRING, fmt.Errorf("Assertion Error.")
	}
	nickName, ok := claims["nne"].(string)
	if !ok {
		return EMPTY_STRING, EMPTY_STRING, EMPTY_STRING, fmt.Errorf("Assertion Error.")
	}

	return userId, maxTimes, nickName, nil
}

/*
	说明：检查短信验证码是否合法
	入参：
	出参：参数1：Token
		 参数1：Error
*/

func CheckSmsCode(phoneNo string, smsCode string) error {
	PrintHead("CheckSmsCode")
	var search smscode.Search
	search.Phone = phoneNo
	search.Status = SMS_STATUS_INIT
	search.SmsCode = smsCode
	r := smscode.New(dbcomm.GetDB(), smscode.DEBUG)
	l, _ := r.GetList(search)
	for _, v := range l {
		local, _ := time.LoadLocation("Local")
		fmt.Println(v.ValidEtime)

		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", v.ValidEtime, local)
		fmt.Println(endTime)
		if endTime.Before(time.Now()) {
			log.Println("过期的，跳过")
			continue
		}
		if v.SmsCode == smsCode {
			m := map[string]interface{}{"status": SMS_STATUS_END}
			err := r.UpdateMap(fmt.Sprintf("%d", v.Id), m, nil)
			if err == nil {
				fmt.Println("OK")
				return nil
			}
		}
		log.Println("非法的。。。。")

	}
	PrintTail("CheckSmsCode")
	return errors.New("短信验证码不合法")

}

/*
	说明：检查短信验证码是否合法
	入参：
	出参：参数1：Token
		 参数1：Error
*/

func CheckCaptchaCode(idKey string, verifyValue string) bool {
	PrintHead("CheckCaptchaCode", idKey, verifyValue)
	verifyResult := base64Captcha.VerifyCaptcha(idKey, verifyValue)
	PrintTail("CheckCaptchaCode", verifyResult)
	return verifyResult
}

/*
	说明：生产二维码信息
	入参：
	出参：参数1：返回符合条件的对象列表
*/

func CreateQrCode(prePayid string, codeUrl string) (string, error) {
	fmt.Println(codeUrl)
	qrCode, err := qr.Encode(codeUrl, qr.M, qr.Auto)
	if err != nil {
		fmt.Println(err)
	}
	// Scale the barcode to 200x200 pixels
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		fmt.Println(err)
	}
	// create the output file
	file, _ := os.Create(prePayid + ".png")
	defer file.Close()
	// encode the barcode as png
	png.Encode(file, qrCode)
	fileBuf, err := ioutil.ReadFile(prePayid + ".png")
	if err != nil {
		fmt.Println(err)
		return EMPTY_STRING, err
	}
	b := bytes.NewBuffer(make([]byte, 0))
	encoder := base64.NewEncoder(base64.StdEncoding, b)
	encoder.Write(fileBuf)
	encoder.Close()
	os.Remove(prePayid + ".png")
	return fmt.Sprintf("data:image/png;base64,%s", b), nil
}

/*
	访问短信平台的请求
*/
type SmsReq struct {
	Data        map[string]string `json:"data"`
	MobliePhone string            `json:"mobilePhone"`
	ServiceType string            `json:"serviceType"`
	Sige        string            `json:"sige"`
	Sys         string            `json:"sys"`
}

/*
	访问短信平台的应答
*/
type SmsResp struct {
	ErrMsg      string `json:"errMsg"`
	Result      string `json:"result"`
	ServiceType string `json:"serviceType"`
}

/*
	说明：调用第三方短信平台
	入参：
	出参：参数1：Token
		 参数1：Error
*/

func PostSmsCode(phoneNo string, smsCode string) error {
	PrintHead("PostSmsCode")
	var smsReq SmsReq
	smsReq.Data = make(map[string]string, 0)
	smsReq.MobliePhone = phoneNo
	smsReq.ServiceType = SMS_SERVICE_TYPE
	smsReq.Sys = SMS_SYS_NO
	smsReq.Data["code"] = smsCode

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(smsReq.Sys + smsReq.ServiceType + smsReq.MobliePhone + SMS_POST_KEY))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	smsReq.Sige = upperSign
	smsReqBuf, err := json.Marshal(smsReq)

	req, err := http.NewRequest("POST", CRF_SMSURL_UAT, bytes.NewReader(smsReqBuf))
	if err != nil {
		log.Println("New Http Request发生错误，原因:", err)
		return errors.New("网络连接出现问题,稍后再试！")
	}
	c := http.Client{}
	resp, _err := c.Do(req)
	if _err != nil {
		log.Println("请求短信平台发送短信错误, 原因:", _err)
		return errors.New("网络连接出现问题,稍后再试！")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return errors.New("获取应答出现问题")
	}
	var smsResp SmsResp
	err = json.Unmarshal(body, &smsResp)
	if err != nil {
		log.Println(err)
	}
	log.Printf("短信发送结果:%#v", smsResp)
	if smsResp.Result != "0000" {
		return errors.New(smsResp.ErrMsg)
	}
	PrintTail("PostSmsCode")
	return nil
}
