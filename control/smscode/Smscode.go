package smscode

import (
	"encoding/json"
	"jcd/control/common"

	"fmt"
	"jcd/service/dbcomm"

	"jcd/service/smscode"
	"math/rand"
	"net/http"
	"time"

	"github.com/mojocn/base64Captcha"
)

/*
	修改密码
*/
type CheckSmsCodeReq struct {
	UserId  string `json:"user_id"`
	SmsCode string `json:"sms_code"`
}

/*
	获取短信验证码--需要图形验证码
*/
type GetSmsCodeReq struct {
	UserId     string `json:"user_id"`
	Phone      string `json:"phone"`
	VerifyCode string `json:"verify_code"`
	IdKey      string `json:"id_key"`
}

/*
	获取短信验证码--需要图形验证码
*/
type GetSmsCodeReqEx struct {
	Phone string `json:"phone"`
}

/*
 */
type GetSmsCodeResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	SmsCode string `json:"sms_code"`
}

/*
	获取图形验证码
*/
type CaptchasReq struct {
	UserId string `json:"user_id"`
}

type CaptchasResp struct {
	ErrCode   string `json:"err_code"`
	ErrMsg    string `json:"err_msg"`
	IdKey     string `json:"id_key"`
	PicBase64 string `json:"pic_base64"`
}

/*
	验证短信验证码
*/
func CheckSmsCode(w http.ResponseWriter, req *http.Request) {
	userId, _, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	fmt.Println("userid=====>", userId)
	var smsCodeReq CheckSmsCodeReq
	var smsCodeResp GetSmsCodeResp
	err := json.NewDecoder(req.Body).Decode(&smsCodeReq)
	if err != nil {
		smsCodeResp.ErrCode = common.ERR_CODE_JSONERR
		smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(smsCodeResp, w, req)
		return
	}
	defer req.Body.Close()

	common.CheckSmsCode("22", smsCodeReq.SmsCode)

}

/*
	生产一个随机的短信验证码
*/
func GetSmsCode(w http.ResponseWriter, req *http.Request) {
	var smsCodeReq GetSmsCodeReq
	var smsCodeResp GetSmsCodeResp
	err := json.NewDecoder(req.Body).Decode(&smsCodeReq)
	if err != nil {
		smsCodeResp.ErrCode = common.ERR_CODE_JSONERR
		smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(smsCodeResp, w, req)
		return
	}

	defer req.Body.Close()

	//验证图形验证码
	//	if !common.CheckCaptchaCode(smsCodeReq.IdKey, smsCodeReq.VerifyCode) {
	//		smsCodeResp.ErrCode = common.ERR_CODE_VERIFY
	//		smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_VERIFY]
	//		common.Write_Response(smsCodeResp, w, req)
	//		return
	//	}

	var search smscode.Search
	search.Phone = smsCodeReq.Phone

	r := smscode.New(dbcomm.GetDB(), smscode.DEBUG)
	//检查上一次发送时间
	l, err := r.GetLast(search)
	if err == nil {
		local, _ := time.LoadLocation("Local")
		intime, _ := time.ParseInLocation("2006-01-02 15:04:05", l.InsertTime, local)
		dt := time.Now().Sub(intime)
		if dt.Seconds() < common.SMSCODE_MIN_INTERVAL {
			smsCodeResp.ErrCode = common.ERR_CODE_TOOBUSY
			smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_TOOBUSY]
			common.Write_Response(smsCodeResp, w, req)
			return
		}
	}
	//增加调用第三方接口
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	err = common.PostSmsCode(smsCodeReq.Phone, vcode)
	if err != nil {
		smsCodeResp.ErrCode = common.ERR_SMS_SENDERR
		smsCodeResp.ErrMsg = err.Error()
		common.Write_Response(smsCodeResp, w, req)
		return
	}
	var e smscode.Smscode
	e.Phone = smsCodeReq.Phone
	e.SmsCode = vcode
	e.InsertTime = time.Now().Format("2006-01-02 15:04:05")
	e.Status = common.SMS_STATUS_INIT
	e.SmsType = common.SMSTYPE_LOGIN
	e.ValidBtime = time.Now().Add(time.Duration(-time.Minute * common.SMSCODE_EXPIRED_MINUTE)).Format("2006-01-02 15:04:05")
	e.ValidEtime = time.Now().Add(time.Duration(time.Minute * common.SMSCODE_EXPIRED_MINUTE)).Format("2006-01-02 15:04:05")
	if err := r.InsertEntity(e, nil); err != nil {
		smsCodeResp.ErrCode = common.ERR_CODE_DBERROR
		smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(smsCodeResp, w, req)
		return
	}
	smsCodeResp.ErrCode = common.ERR_CODE_SUCCESS
	smsCodeResp.SmsCode = vcode
	smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS] + "获取短信验证码成功"
	common.Write_Response(smsCodeResp, w, req)
}

/*
	生产一个随机的短信验证码

*/
func GetSmsCodeEx(w http.ResponseWriter, req *http.Request) {

	var smsCodeReq GetSmsCodeReqEx
	var smsCodeResp GetSmsCodeResp
	err := json.NewDecoder(req.Body).Decode(&smsCodeReq)
	if err != nil {
		smsCodeResp.ErrCode = common.ERR_CODE_JSONERR
		smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(smsCodeResp, w, req)
		return
	}

	defer req.Body.Close()
	var search smscode.Search
	if smsCodeReq.Phone == common.EMPTY_STRING {
		smsCodeResp.ErrCode = common.ERR_CODE_JSONERR
		smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR]
		common.Write_Response(smsCodeResp, w, req)
		return
	}
	search.Phone = smsCodeReq.Phone
	r := smscode.New(dbcomm.GetDB(), smscode.DEBUG)
	//检查上一次发送时间
	l, err := r.GetLast(search)
	if err == nil {
		local, _ := time.LoadLocation("Local")
		intime, _ := time.ParseInLocation("2006-01-02 15:04:05", l.InsertTime, local)
		dt := time.Now().Sub(intime)
		if dt.Seconds() < common.SMSCODE_MIN_INTERVAL {
			smsCodeResp.ErrCode = common.ERR_CODE_TOOBUSY
			smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_TOOBUSY]
			common.Write_Response(smsCodeResp, w, req)
			return
		}
	}
	//增加调用第三方接口
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	err = common.PostSmsCode(smsCodeReq.Phone, vcode)
	if err != nil {
		smsCodeResp.ErrCode = common.ERR_SMS_SENDERR
		smsCodeResp.ErrMsg = err.Error()
		common.Write_Response(smsCodeResp, w, req)
		return
	}
	var e smscode.Smscode
	e.Phone = smsCodeReq.Phone
	e.SmsCode = vcode
	e.InsertTime = time.Now().Format("2006-01-02 15:04:05")
	e.Status = common.SMS_STATUS_INIT
	e.SmsType = common.SMSTYPE_LOGIN
	e.ValidBtime = time.Now().Add(time.Duration(-time.Minute * common.SMSCODE_EXPIRED_MINUTE)).Format("2006-01-02 15:04:05")
	e.ValidEtime = time.Now().Add(time.Duration(time.Minute * common.SMSCODE_EXPIRED_MINUTE)).Format("2006-01-02 15:04:05")
	if err := r.InsertEntity(e, nil); err != nil {
		smsCodeResp.ErrCode = common.ERR_CODE_DBERROR
		smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(smsCodeResp, w, req)
		return
	}
	smsCodeResp.ErrCode = common.ERR_CODE_SUCCESS
	smsCodeResp.SmsCode = vcode
	smsCodeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS] + "获取短信验证码成功"
	common.Write_Response(smsCodeResp, w, req)
}

/*
	获取图形验证码
*/
func GetCaptchas(w http.ResponseWriter, r *http.Request) {
	var req CaptchasReq
	var resp CaptchasResp
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.ErrCode = common.ERR_CODE_JSONERR
		resp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(resp, w, r)
		return
	}
	defer r.Body.Close()

	//数字验证码配置
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      200,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}

	idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	//以base64编码
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)
	resp.PicBase64 = base64stringD
	resp.IdKey = idKeyD
	resp.ErrCode = common.ERR_CODE_SUCCESS
	resp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS] + "获取图形验证码成功"
	common.Write_Response(resp, w, r)
}
