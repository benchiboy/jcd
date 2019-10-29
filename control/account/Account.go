package account

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"jcd/control/common"
	"jcd/service/account"
	"jcd/service/dbcomm"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
	查询账户请求
*/
type GetAccountReq struct {
	LoginName string `json:"login_name"`
}

/*
	查询账户返回
*/
type GetAccountResp struct {
	ErrCode     string  `json:"err_code"`
	ErrMsg      string  `json:"err_msg"`
	NickName    string  `json:"nick_name"`
	AvatarUrl   string  `json:"avatar_url"`
	Mail        string  `json:"mail"`
	Phone       string  `json:"phone"`
	Gender      int     `json:"gender"`
	AccountBal  float64 `json:"account_bal"`
	Country     string  `json:"country"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	PassProblem string  `json:"pass_problem"`
}

/*
	查询账户是否存在的请求
*/
type CheckAccountReq struct {
	LoginName string `json:"login_name"`
}

/*
	查询账户是否存在的请求返回
*/
type CheckAccountResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
	查询账户是否存在的请求返回
*/
type CheckSignInResp struct {
	ErrCode  string `json:"err_code"`
	ErrMsg   string `json:"err_msg"`
	NickName string `json:"nick_name"`
	userId   string `json:"user_id"`
}

/*
	查询账户是否存在的请求返回
*/
type GetAvatarUrlResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
	注册请求
*/
type SignUpReq struct {
	NickName    string `json:"nick_name"`
	UserName    string `json:"login_name"`
	PassWord    string `json:"login_pass"`
	SmsCode     string `json:"sms_code"`
	PassProblem string `json:"pass_problem"`
	PassAnswer  string `json:"pass_answer"`
	HeadImage   string `json:"head_image"`
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

func CheckSignIn(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("CheckSignIn")
	userId, _, nickName, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	var signInResp CheckSignInResp
	signInResp.NickName = nickName
	signInResp.userId = userId
	signInResp.ErrCode = common.ERR_CODE_SUCCESS
	signInResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(signInResp, w, req)
	common.PrintTail("CheckSignIn")
}

/*
	说明： 检查注册用户是否存在
	出参： 返回用户的登录状态
*/

func GetAccount(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("GetAccount")
	var accountReq CheckAccountReq
	var accountResp CheckAccountResp
	err := json.NewDecoder(req.Body).Decode(&accountReq)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer req.Body.Close()
	var search account.Search
	search.LoginName = accountReq.LoginName
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	if _, err := r.Get(search); err == nil {
		accountResp.ErrCode = common.ERR_CODE_EXISTED
		accountResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED]
		common.Write_Response(accountResp, w, req)
		return
	}
	//需要用户注册
	accountResp.ErrCode = common.CODE_SUCC
	accountResp.ErrMsg = common.ERROR_MAP[common.CODE_SUCC]
	common.Write_Response(accountResp, w, req)
	common.PrintTail("GetAccount")
}

/*
	说明： 得到登录账户信息
	出参： 返回用户的登录状态
*/

func GetAccountInfo(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("GetAccountInfo")
	//	userId, _, _, tokenErr := common.CheckToken(w, req)
	//	if tokenErr != nil {
	//		return
	//	}
	//	uId, _ := strconv.ParseInt(userId, 10, 64)

	var accountReq GetAccountReq
	var accountResp GetAccountResp
	err := json.NewDecoder(req.Body).Decode(&accountReq)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer req.Body.Close()
	log.Println(accountReq)
	var search account.Search
	search.LoginName = accountReq.LoginName
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	if e, err := r.Get(search); err == nil {
		accountResp.ErrCode = common.CODE_SUCC
		accountResp.ErrMsg = common.ERROR_MAP[common.CODE_SUCC]
		accountResp.NickName = e.NickName
		accountResp.AvatarUrl = e.AvatarUrl
		accountResp.AccountBal = e.AccountBal
		accountResp.Gender = e.Gender
		accountResp.Mail = e.Mail
		accountResp.Phone = e.Phone
		accountResp.Country = e.Country
		accountResp.Province = e.Province
		accountResp.City = e.City
		accountResp.PassProblem = e.PassProblem

		common.Write_Response(accountResp, w, req)
		return
	}

	common.PrintTail("GetAccountInfo")
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

	var search account.Search
	search.LoginName = signupReq.UserName
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	if e, err := r.Get(search); err == nil {
		signupResp.ErrCode = common.ERR_CODE_EXISTED
		signupResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED] + e.LoginName
		common.Write_Response(signupResp, w, req)
		return

	} else {
		//		if err := common.CheckSmsCode(signupReq.UserName, signupReq.SmsCode); err != nil {
		//			signupResp.ErrCode = common.ERR_CODE_VERIFY
		//			signupResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_VERIFY]
		//			common.Write_Response(signupResp, w, req)
		//			return
		//		}
		var e account.Account
		e.LoginMode = common.LOGIN_PHONE
		e.LoginName = signupReq.UserName
		e.LoginPass = signupReq.PassWord
		e.PassProblem = signupReq.PassProblem
		e.PassAnswer = signupReq.PassAnswer
		e.AvatarUrl = signupReq.HeadImage
		e.UserId = time.Now().Unix()
		e.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
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
	说明：根据网页的提交，修改账户信息
	出参：
*/
func SetAccount(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("SetAccount")
	userId, _, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	var form account.Account
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println(form)
	defer req.Body.Close()
	var search account.Search
	search.UserId = uId
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	if e, err := r.Get(search); err == nil {
		r.UpdataEntity(fmt.Sprintf("%d", e.Id), form, nil)
	}
	common.Write_Response("OK", w, req)
	common.PrintTail("SetAccount")
}

/*
	说明：根据网页的提交，修改账户信息
	出参：
*/
func GetAvatarUrl(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("GetAvatarUrl")
	var avatarResp GetAvatarUrlResp
	avatarId, ok := req.URL.Query()["id"]
	if !ok || len(avatarId) < 1 {
		avatarResp.ErrCode = common.ERR_CODE_URLERR
		avatarResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_URLERR] + "Get b_id错误！"
		common.Write_Response(avatarResp, w, req)
		return
	}
	imagBuf, _ := ioutil.ReadFile("./images/" + avatarId[0] + ".jpeg")
	w.Write(imagBuf)
	common.PrintTail("GetAvatarUrl")
}

/*
	1、上传用户头像图片
	2、如果上传成功，更新DB中存储文件名称
*/
func UploadPics(w http.ResponseWriter, r *http.Request) {
	//	log.Println(comm.BEGIN_TAG, "UploadPics......")
	r.FormFile("file")
	file, handle, err := r.FormFile("file")
	if err != nil {
		err.Error()
	}
	fmt.Println(handle.Filename, handle.Size)
	f, err := os.OpenFile("22222", os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		err.Error()
	}
	defer f.Close()
	defer file.Close()
	fmt.Println("upload success")

	//log.Println(comm.END_TAG, "UploadPics Successful......", time.Since(t1))
	return
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
