package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

const STATUS_DISABLED = 1
const STATUS_ENABLED = 0
const MAX_SEARCH_TIMES = "5"
const STATUS_INIT = "i"

const FIELD_LOGIN_PASS = "login_pass"
const FIELD_ERRORS = "errors"

const DEFAULT_PWD = "123456"

const EMPTY_STRING = ""

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
