package login

import (
	"encoding/json"
	"fmt"
	"jcd/control/common"
	"jcd/service/account"

	"jcd/service/dbcomm"
	"jcd/service/login"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
 */
type Login struct {
	Code string `json:"code"`
}

/*
 */
type LoginResp struct {
	Openid      string `json:"openid"`
	Session_key string `json:"session_key"`
	Unionid     string `json:"unionid"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

/*
	根据CODE 得到OPENID
*/

func wxGetOpenid(code string) (error, string, string, string) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	getUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=wxcc7ef55685a5221c&secret=4d53e212c52cd1955703cf45600f7472&js_code=" + code + "&grant_type=authorization_code"
	res, err := httpClient.Get(getUrl)
	if err != nil {
		return fmt.Errorf("访问微信认证服务出错！"), "", "", ""
	}
	defer res.Body.Close()
	var resp LoginResp
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return fmt.Errorf("解析JSON出错"), "", "", ""
	}
	log.Printf("%#v", resp)
	return nil, resp.Openid, resp.Unionid, resp.Session_key
}

/*
	微信登录
*/
func WxLogin(w http.ResponseWriter, req *http.Request) {
	log.Println("========》WxLogin")
	keys, ok := req.URL.Query()["code"]
	if !ok || len(keys) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	code := keys[0]
	err, openId, unionId, sessionKey := wxGetOpenid(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uId, err := regUser(openId, unionId, sessionKey)
	w.Write([]byte(uId))
	log.Println("《========WxLogin")
}

/*
	登记用户注册信息
*/
func regUser(openId string, unionId string, sessionKey string) (string, error) {
	var search account.Search
	search.PuserId = openId
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	var userId string
	if e, err := r.Get(search); err != nil {
		var a account.Account
		a.PuserId = openId
		a.PunionId = unionId
		a.UserId = time.Now().Unix()
		a.PsessionKey = sessionKey
		a.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
		r.InsertEntity(a, nil)
		userId = fmt.Sprintf("%d", a.UserId)
	} else {
		r := login.New(dbcomm.GetDB(), login.DEBUG)
		var l login.Login
		l.UserId = e.UserId
		l.LoginTime = time.Now().Format("2006-01-02 15:04:05")
		l.LoginNo = time.Now().Unix()
		r.InsertEntity(l, nil)
		userId = fmt.Sprintf("%d", e.UserId)
	}
	return userId, nil
}

/*
 */
type SigninReq struct {
	UserName string `json:"login_name"`
	PassWord string `json:"login_pass"`
}

/*
 */
type SigninResp struct {
	ErrCode       string `json:"err_code"`
	ErrMsg        string `json:"err_msg"`
	Token         string `json:"token"`
	NickName      string `json:"nick_name"`
	UserId        int64  `json:"user_id"`
	LastLoginTime string `json:"last_login_time"`
}

/*
 */
type CheckSigninResp struct {
	UserId   int64  `json:"user_id"`
	NickName string `json:"nick_name"`
	IsSignIn bool   `json:"is_signin"`
	ErrCode  string `json:"err_code"`
	ErrMsg   string `json:"err_msg"`
}

/*
 */
type LogoutResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
	检查是否登录
*/
func CheckSignIn(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("CheckSignIn")
	userId, _, nickName, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	var checkSigninResp CheckSigninResp
	checkSigninResp.NickName = nickName
	checkSigninResp.UserId = uId
	checkSigninResp.IsSignIn = true

	checkSigninResp.ErrCode = common.CODE_SUCC
	checkSigninResp.ErrMsg = common.ERROR_MAP[common.CODE_SUCC]
	common.Write_Response(checkSigninResp, w, req)

	common.PrintHead("CheckSignIn")
}

/*
	系统登录
*/
func SignIn(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("SignIn")
	var login SigninReq
	var loginResp SigninResp
	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		loginResp.ErrCode = common.ERR_CODE_JSONERR
		loginResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(loginResp, w, req)
		return
	}
	defer req.Body.Close()
	var search account.Search
	search.LoginName = login.UserName
	a := account.New(dbcomm.GetDB(), account.DEBUG)
	e, err := a.Get(search)
	if err != nil || e == nil {
		loginResp.ErrCode = common.ERR_CODE_NOTFIND
		loginResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND] + " 账户不存在"
		common.Write_Response(loginResp, w, req)
		return
	}

	if e.Status == common.STATUS_DISABLED {
		loginResp.ErrCode = common.ERR_CODE_STATUS
		loginResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_STATUS] + "账户已经被冻结！"
		common.Write_Response(loginResp, w, req)
		return
	}

	if e.LoginPass != login.PassWord {
		loginResp.ErrCode = common.ERR_CODE_NOMATCH
		loginResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOMATCH] + login.UserName + "用户名或密码错误"
		common.Write_Response(loginResp, w, req)
		return
	}

	tokenStr, err := common.GetToken(fmt.Sprintf("%d", e.UserId), e.NickName)
	if err != nil {
		loginResp.ErrCode = common.ERR_CODE_TOKENER
		loginResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_TOKENER] + err.Error()
		common.Write_Response(loginResp, w, req)
		return
	}

	loginResp.ErrCode = common.ERR_CODE_SUCCESS
	loginResp.NickName = e.NickName
	loginResp.UserId = e.UserId
	loginResp.ErrMsg = "用户 " + login.UserName + " 登录成功！"
	loginResp.Token = tokenStr
	common.PrintTail("SignIn")
	common.Write_Response(loginResp, w, req)
}

/*
	系统退出

*/
func SignOut(w http.ResponseWriter, req *http.Request) {
	_, _, _, err := common.CheckToken(w, req)
	if err != nil {
		return
	}

	var resp LogoutResp
	resp.ErrCode = common.ERR_CODE_SUCCESS
	resp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS] + "退出成功！"
	common.Write_Response(resp, w, req)

}
