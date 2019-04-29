package account

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"jcd/control/common"
	"jcd/service/account"
	"jcd/service/dbcomm"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
	查询账户请求
*/
type GetAccountReq struct {
	UserId int64 `json:"user_id"`
}

/*
	查询账户返回
*/
type GetAccountResp struct {
	ErrCode   string `json:"err_code"`
	ErrMsg    string `json:"err_msg"`
	NickName  string `json:"nick_name"`
	AvatarUrl string `json:"avatar_url"`
}

/*
	注册请求
*/
type SignUpReq struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
	SmsCode  string `json:"sms_code"`
}

/*
	注册返回
*/
type SignUpResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type EncryptedDataUserInfo struct {
	OpenID    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	UnionID   string `json:"unionId"`
	Watermark struct {
		Timestamp int    `json:"timestamp"`
		Appid     string `json:"appid"`
	} `json:"watermark"`
}

/*
	说明：检查用户是否登录
	出参： 返回用户的登录状态
*/

func GetAccount(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("GetAccount")
	var accountReq GetAccountReq
	var accountResp GetAccountResp
	err := json.NewDecoder(req.Body).Decode(&accountReq)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer req.Body.Close()
	var search account.Search
	search.UserId = accountReq.UserId
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	if e, err := r.Get(search); err == nil {
		//用户存在，已经登录
		if e.ExpiresIn > time.Now().Unix() {
			accountResp.ErrCode = common.ERR_USER_SIGNINED
			accountResp.ErrMsg = common.ERROR_MAP[common.ERR_USER_SIGNINED]
			accountResp.NickName = e.NickName
			accountResp.AvatarUrl = e.AvatarUrl
			common.Write_Response(accountResp, w, req)
			return
			//用户存在，未登录
		} else {
			accountResp.ErrCode = common.ERR_USER_UNSIGNIN
			accountResp.ErrMsg = common.ERROR_MAP[common.ERR_USER_UNSIGNIN] + e.LoginName
			common.Write_Response(accountResp, w, req)
			return
		}
	}
	//需要用户注册
	accountResp.ErrCode = common.ERR_USER_MSTSIGNUP
	accountResp.ErrMsg = common.ERROR_MAP[common.ERR_USER_MSTSIGNUP]
	common.Write_Response(accountResp, w, req)
	common.PrintTail("GetAccount")
}

/*
	说明：账户注册
	出参：参数1：返回符合条件的对象列表
*/

func SignUp(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("SignUp")
	var signupReq SignUpReq
	var signupResp SignUpResp
	err := json.NewDecoder(req.Body).Decode(&signupReq)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer req.Body.Close()
	//	if err := common.CheckSmsCode(0, signupReq.UserName, signupReq.SmsCode); err != nil {
	//		signupResp.ErrCode = common.ERR_CODE_VERIFY
	//		signupResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_VERIFY]
	//		common.Write_Response(signupResp, w, req)
	//		return
	//	}
	var search account.Search
	search.LoginName = signupReq.UserName
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	if e, err := r.Get(search); err == nil {
		signupResp.ErrCode = common.ERR_CODE_EXISTED
		signupResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED] + e.LoginName
		common.Write_Response(signupResp, w, req)
		return

	} else {
		var e account.Account
		e.LoginName = signupReq.UserName
		e.LoginPass = signupReq.PassWord
		e.UserId = time.Now().Unix()
		r.InsertEntity(e, nil)
	}
	signupResp.ErrCode = common.ERR_CODE_SUCCESS
	signupResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(signupResp, w, req)

	common.PrintTail("SignUp")
}

/*

 */

func OAuthSignUp(inAaccount account.Account) {
	common.PrintHead("OAuthSignUp")
	var search account.Search
	search.PuserId = inAaccount.PuserId
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	if e, err := r.Get(search); err == nil {
		fmt.Println("此用户已经注册====>", e)
	} else {
		inAaccount.UserId = time.Now().Unix()
		inAaccount.LoginMode = common.LOGIN_OAUTH
		r.InsertEntity(inAaccount, nil)
	}
	common.PrintTail("OAuthSignUp")
}

/*
	说明：更新账号信息
	出参：参数1：返回符合条件的对象列表
*/

func UpdateAccount(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("UpdateAccount")
	var form account.Form
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer req.Body.Close()
	var search account.Search
	search.UserId = form.Form.UserId
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	if e, err := r.Get(search); err == nil {
		u := getWechatUserInfo(form.Form.EncryptedData, form.Form.Iv, e.PsessionKey)
		e.AvatarUrl = u.AvatarURL
		e.Province = u.Province
		e.City = u.City
		e.Country = u.Country
		e.Language = u.Language
		e.Gender = u.Gender
		r.UpdataEntity(fmt.Sprintf("%d", e.Id), *e, nil)
	} else {
		r.InsertEntity(form.Form, nil)
	}
	common.Write_Response("OK", w, req)
	common.PrintTail("UpdateAccount")
}

/*
	说明：得到微信的基本信息
	出参：参数1：返回符合条件的对象列表
*/

func getWechatUserInfo(inEncryptedData string, inIv string, inSessionKey string) *EncryptedDataUserInfo {
	common.PrintHead("getWechatUserInfo")
	encryptedData, _ := base64.StdEncoding.DecodeString(inEncryptedData)
	iv, _ := base64.StdEncoding.DecodeString(inIv)
	sessionKey, _ := base64.StdEncoding.DecodeString(inSessionKey)

	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err := aes.NewCipher([]byte(sessionKey))
	if err != nil {
		return nil
	}
	decrypted := make([]byte, len(encryptedData))
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.CryptBlocks(decrypted, encryptedData)
	var userInfo EncryptedDataUserInfo
	t := string(decrypted)
	fmt.Println(t)
	total := strings.Index(t, "}}") + 2
	err = json.Unmarshal(decrypted[:total], &userInfo)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	log.Println(userInfo.OpenID)
	common.PrintTail("getWechatUserInfo")
	return &userInfo

}
