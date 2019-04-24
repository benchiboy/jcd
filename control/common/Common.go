package common

import (
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
)

var (
	WX_PAY_URL          = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	WX_PAY_CALLBACK_URL = "http://132.232.11.85:8087/jc/api/wxcallback"
	MCT_ID              = "1452819402"
	APP_ID              = "wx2db791be2eb77467"
	PRODUCT_NAME        = "测试商品"
	SERVER_IP           = "132.232.11.85"
	MCT_KEY             = "qj837vwk83xk2902jidk93slw82ms5ka"
	TRADE_TYPE_NATIVE   = "NATIVE"
	TRADE_TYPE_JSAPI    = "JSAPI"
)

const ERR_CODE_SUCCESS = "0000"
const ERR_CODE_DBERROR = "1001"
const ERR_CODE_TOKENER = "1003"
const ERR_CODE_PARTOEN = "1005"
const ERR_CODE_JSONERR = "2001"
const ERR_CODE_URLERR = "2005"
const ERR_CODE_NOTFIND = "3000"
const ERR_CODE_NOMATCH = "3010"
const ERR_CODE_EXPIRED = "8000"
const ERR_CODE_TYPEERR = "4000"
const ERR_CODE_STATUS = "5000"
const ERR_CODE_FAILED = "9000"
const ERR_CODE_OPERTYP = "4005"
const ERR_CODE_EXISTED = "4040"
const ERR_CODE_TOOBUSY = "6010"
const ERR_CODE_VERIFY = "7020"

const STATUS_DISABLED = 1
const STATUS_ENABLED = 0
const STATUS_SUCC = "S"
const STATUS_FAIL = "F"
const MAX_SEARCH_TIMES = "5"
const STATUS_INIT = "i"

const FIELD_LOGIN_PASS = "login_pass"
const FIELD_ERRORS = "errors"
const FIELD_KILLS = "kills"
const FIELD_LIKES = "likes"
const FIELD_UPDATE_TIME = "update_time"
const FIELD_UPDATE_USER = "update_user"
const FIELD_PROC_STATUS = "proc_status"
const FIELD_PROC_MSG = "proc_msg"

const DEFAULT_PWD = "123456"
const SUCC_MSG = "success"
const EMPTY_STRING = ""

const SMSTYPE_LOGIN = "login"
const SMSTYPE_RESET = "reset"
const SMS_STATUS_INIT = "i"
const SMS_STATUS_END = "e"
const SMSCODE_EXPIRED_MINUTE = 5
const SMSCODE_MIN_INTERVAL = 10

const COMMENT_INIT_VALUE = 0
const COMMENT_LIKE = 10
const COMMENT_KILL = 20
const COMMENT_REPLY = 30

var (
	ERROR_MAP map[string]string = map[string]string{
		ERR_CODE_SUCCESS: "执行成功:",
		ERR_CODE_DBERROR: "DB执行错误:",
		ERR_CODE_JSONERR: "JSON格式错误:",
		ERR_CODE_EXPIRED: "时效已经到期:",
		ERR_CODE_TYPEERR: "类型转换错误:",
		ERR_CODE_STATUS:  "状态不正确:",
		ERR_CODE_TOKENER: "获取TOKEN失败:",
		ERR_CODE_PARTOEN: "解析TOKEN错误:",
		ERR_CODE_NOMATCH: "比较不匹配:",
		ERR_CODE_URLERR:  "Url传参有误:",
		ERR_CODE_OPERTYP: "ShowType类型错误:",
		ERR_CODE_NOTFIND: "查询没发现提示:",
		ERR_CODE_EXISTED: "注册账户已经存在:",
		ERR_CODE_TOOBUSY: "短信发送太频繁:",
		ERR_CODE_VERIFY:  "验证码校验错误:",
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
	w.Header().Set("Access-Control-Allow-Origin", "http://10.89.4.225:8000")
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

func GetToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["aud"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(10)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["cnt"] = MAX_SEARCH_TIMES

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

func CheckToken(w http.ResponseWriter, req *http.Request) (string, string, error) {
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
		return EMPTY_STRING, EMPTY_STRING, err
	}
	userId, ok := claims["aud"].(string)
	if !ok {
		errResp.ErrCode = ERR_CODE_TYPEERR
		errResp.ErrMsg = ERROR_MAP[ERR_CODE_TYPEERR] + "userid"
		Write_Response(errResp, w, req)

		return EMPTY_STRING, EMPTY_STRING, fmt.Errorf("Assertion Error.")
	}
	maxTimes, ok := claims["cnt"].(string)
	if !ok {
		errResp.ErrCode = ERR_CODE_TYPEERR
		errResp.ErrMsg = ERROR_MAP[ERR_CODE_TYPEERR] + "maxTimes"
		Write_Response(errResp, w, req)
		return EMPTY_STRING, EMPTY_STRING, fmt.Errorf("Assertion Error.")
	}

	return userId, maxTimes, nil
}

/*
	说明：检查短信验证码是否合法
	入参：
	出参：参数1：Token
		 参数1：Error
*/

func CheckSmsCode(userId int64, phone string, smsCode string) error {
	PrintHead("CheckSmsCode")
	var search smscode.Search
	search.Phone = phone
	search.UserId = userId
	search.Status = SMS_STATUS_INIT
	search.SmsCode = smsCode
	r := smscode.New(dbcomm.GetDB(), smscode.DEBUG)
	l, _ := r.GetList(search)
	for _, v := range l {
		local, _ := time.LoadLocation("Local")
		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", v.ValidEtime, local)
		if endTime.Before(time.Now()) {
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
	PrintHead("CheckCaptchaCode")
	verifyResult := base64Captcha.VerifyCaptcha(idKey, verifyValue)
	PrintTail("CheckCaptchaCode")
	return verifyResult
}
