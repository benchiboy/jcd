package oauth

import (
	"encoding/json"
	"fmt"
	"jcd/control/account"
	"jcd/control/common"
	account_service "jcd/service/account"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

/*
 */
type WeiboUserResp struct {
	Id                int64  `json:"id"`
	Idstr             string `json:"idstr"`
	Screen_name       string `json:"screen_name"`
	Name              string `json:"name"`
	Province          string `json:"province"`
	City              string `json:"city"`
	Location          string `json:"location"`
	Description       string `json:"description"`
	Profile_image_url string `json:"profile_image_url"`
	Gender            string `json:"gender"`
	Lang              string `json:"lang"`
}

var endpotin = oauth2.Endpoint{
	AuthURL:  "https://api.weibo.com/oauth2/authorize",
	TokenURL: "https://api.weibo.com/oauth2/access_token",
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     "2302557195",
	ClientSecret: "3e096effd83cd4d0553b5098159eced8",
	RedirectURL:  "http://www.doulaikan.club/jc/api/weibocallback",
	Scopes:       []string{"all"},
	Endpoint:     endpotin,
}

const oauthStateString = "random"

/*
	微博OAUTH 用户注册
*/
func WeiboSignUp(userResp WeiboUserResp) {
	var e account_service.Account
	e.PuserId = userResp.Idstr
	e.AvatarUrl = userResp.Profile_image_url
	e.NickName = userResp.Screen_name
	e.Language = userResp.Lang
	account.OAuthSignUp(e)
}

/*
	微信OAUTH 用户注册
*/
func WxSignUp(userResp WeiboUserResp) {
	var e account_service.Account
	e.PuserId = userResp.Idstr
	e.AvatarUrl = userResp.Profile_image_url
	e.NickName = userResp.Screen_name
	e.Language = userResp.Lang
	account.OAuthSignUp(e)
}

/*
	QQOAUTH 用户注册
*/
func QqSignUp(userResp WeiboUserResp) {
	var e account_service.Account
	e.PuserId = userResp.Idstr
	e.AvatarUrl = userResp.Profile_image_url
	e.NickName = userResp.Screen_name
	e.Language = userResp.Lang
	account.OAuthSignUp(e)
}

/*
	微博OAUTH 的回调
*/
func WeiboCallback(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("WeiboCallback")
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println(token.AccessToken)
	fmt.Println(token.RefreshToken)
	fmt.Println(token.Expiry.Format("2006-01-02 15:04:05"))

	userUrl := "https://api.weibo.com/2/users/show.json?access_token=" + token.AccessToken +
		"&uid=" + token.Extra("uid").(string)
	fmt.Println(userUrl)
	res, err := http.Get(userUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	var userResp WeiboUserResp
	err = json.NewDecoder(res.Body).Decode(&userResp)
	if err != nil {
		fmt.Println(err)
		return
	}
	WeiboSignUp(userResp)

	common.PrintTail("WeiboCallback")

	//	http.Redirect(w, r, "http://www.doulaikan.club/jc/api/home", http.StatusFound)
	return
}

/*
	微博OAUTH 的回调
*/

func WeiboCancelCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("......proc user reg......")
	t1 := time.Now()
	log.Println(t1)
	return
}

/*
	微博OAUTH 到登录页面
*/

func WeiboLogin(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("WeiboLogin")
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusFound)
	common.PrintTail("WeiboLogin")
}

func WeiboLoginIndex(w http.ResponseWriter, r *http.Request) {

	htmlIndex := `<html><body>
<a href="/jc/api/weibologin">Welcome to login WEIBO</a>
</body></html>`
	fmt.Fprintf(w, htmlIndex)

}

///*
//	根据CODE 得到OPENID
//*/
//func wxGetOpenid(code string) (error, string, string, string) {
//	httpClient := &http.Client{
//		Timeout: 10 * time.Second,
//	}
//	getUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=wx384c33f61f535883&secret=c885c58b76fdda0135adf6f42e32b27a&js_code=" + code + "&grant_type=authorization_code"
//	res, err := httpClient.Get(getUrl)
//	if err != nil {
//		return fmt.Errorf("访问微信认证服务出错！"), "", "", ""
//	}
//	defer res.Body.Close()
//	var resp LoginResp
//	err = json.NewDecoder(res.Body).Decode(&resp)
//	if err != nil {
//		return fmt.Errorf("解析JSON出错"), "", "", ""
//	}
//	log.Printf("%#v", resp)
//	return nil, resp.Openid, resp.Unionid, resp.Session_key
//}

///*
//	微信登录
//*/
//func WxLogin(w http.ResponseWriter, req *http.Request) {
//	log.Println("========》WxLogin")
//	keys, ok := req.URL.Query()["code"]
//	if !ok || len(keys) < 1 {
//		log.Println("Url Param 'key' is missing")
//		return
//	}
//	code := keys[0]
//	err, openId, unionId, sessionKey := wxGetOpenid(code)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	uId, err := regUser(openId, unionId, sessionKey)
//	w.Write([]byte(uId))
//	log.Println("《========WxLogin")
//}

func WxLoginIndex(w http.ResponseWriter, r *http.Request) {

	htmlIndex := `<html><body>
		<a href="/jc/api/wxlogin">Welcome to login WEIBO</a>
		</body></html>`
	fmt.Fprintf(w, htmlIndex)

}

func Home(w http.ResponseWriter, r *http.Request) {

	htmlIndex := `<html><body>
		<h1>555555555555555555</h1>
		</body></html>`
	fmt.Fprintf(w, htmlIndex)

}

func WxLogin(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("WxLogin")

	//	goUrl := "https: //open.weixin.qq.com/connect/oauth2/authorize?appid=wx384c33f61f535883&redirect_uri=http://132.232.11.85:8087/jc/api//jc/api/wxcallback&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect"

	common.PrintTail("WxLogin")
}

func WxCallback(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("WxCallback")

	common.PrintTail("WxCallback")
}
