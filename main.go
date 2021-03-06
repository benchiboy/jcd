// rcs_contract_mgr project main.go
package main

import (
	"flag"

	"io"

	"jcd/control/account"
	"jcd/control/badloan"
	"jcd/control/comment"
	"jcd/control/flow"
	"jcd/control/home"
	"jcd/control/index"
	"jcd/control/login"
	"jcd/control/oauth"
	"jcd/control/payutil"
	"jcd/control/pwd"
	"jcd/control/region"
	"jcd/control/tiplist"

	"jcd/control/smscode"

	"jcd/service/dbcomm"
	//"jcd/util"
	"log"
	"net/http"
	"os"

	goconf "github.com/pantsing/goconf"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	http_srv   *http.Server
	dbUrl      string
	ccdbUrl    string
	listenPort int
	idleConns  int
	openConns  int
)

func go_WebServer() {
	log.Println("Listen Service start...")

	//账号相关
	http.HandleFunc("/jc/api/weiboindex", oauth.WeiboLoginIndex)
	http.HandleFunc("/jc/api/weibologin", oauth.WeiboLogin)
	http.HandleFunc("/jc/api/weibocallback", oauth.WeiboCallback)
	http.HandleFunc("/jc/api/weibocancelcallback", oauth.WeiboCancelCallback)

	http.HandleFunc("/jc/api/alipaycallback", oauth.AlipayCallback)
	http.HandleFunc("/jc/api/alipaylogin", oauth.AlipayLogin)

	http.HandleFunc("/wxLogin", login.WxLogin)
	http.HandleFunc("/jc/api/wxlogin", oauth.WeiboLogin)
	http.HandleFunc("/jc/api/wxloginindex", oauth.WeiboLoginIndex)

	http.HandleFunc("/jc/api/getaccountinfo", account.GetAccountInfo)
	http.HandleFunc("/jc/api/getaccount", account.GetAccount)
	http.HandleFunc("/jc/api/setaccount", account.SetAccount)
	http.HandleFunc("/jc/api/getavatar", account.GetAvatarUrl)

	http.HandleFunc("/jc/api/tiplist", tiplist.GetTipList)
	http.HandleFunc("/jc/api/disputes", flow.Disputes)
	http.HandleFunc("/jc/api/signin", login.SignIn)
	http.HandleFunc("/jc/api/checksignin", account.CheckSignIn)

	http.HandleFunc("/jc/api/signout", login.SignOut)
	http.HandleFunc("/jc/api/signup", account.SignUp)
	http.HandleFunc("/jc/api/resetpwd", pwd.ResetPwd)
	http.HandleFunc("/jc/api/chgpwd", pwd.ChangePwd)
	//支付相关
	http.HandleFunc("/jc/api/wxpaycallback", payutil.WxpayCallback)
	http.HandleFunc("/jc/api/findloan", flow.FindLoans)
	http.HandleFunc("/jc/api/repay", flow.RepayOrder)
	http.HandleFunc("/jc/api/cpcnrepay", flow.CPCNRepayOrder)
	http.HandleFunc("/jc/api/myflowlist", flow.MyFlowList)
	http.HandleFunc("/jc/api/myfeedbacklist", flow.MyFeedbackList)
	http.HandleFunc("/jc/api/orderquery", flow.OrderStatus)
	//验证码相关
	http.HandleFunc("/jc/api/getsmscode", smscode.GetSmsCodeEx)
	http.HandleFunc("/jc/api/checksmscode", smscode.CheckSmsCode)
	http.HandleFunc("/jc/api/getcaptcha", smscode.GetCaptchas)
	//评论相关
	http.HandleFunc("/jc/api/mycommlist", comment.MyCommentList)
	http.HandleFunc("/jc/api/commlist", comment.CommentList)
	http.HandleFunc("/jc/api/postcomm", comment.CommentPost)
	http.HandleFunc("/jc/api/likecomm", comment.CommentLike)
	http.HandleFunc("/jc/api/killcomm", comment.CommentKill)
	http.HandleFunc("/jc/api/replycomm", comment.CommentReply)

	//逾期指数接口
	http.HandleFunc("/jc/api/badfact", index.BadPLoanList)
	http.HandleFunc("/jc/api/badploan", badloan.BadPLoanList)
	http.HandleFunc("/jc/api/home", home.Home)
	//基础参数信息
	http.HandleFunc("/jc/api/getregion", region.GetRegionList)

	http_srv = &http.Server{
		Addr: ":8087",
	}
	log.Printf("listen:")
	if err := http_srv.ListenAndServe(); err != nil {
		log.Printf("listen: %s\n", err)
	}
}

func init() {
	log.Println("System Paras Init......")
	log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)
	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "jcd.log",
		MaxSize:    500, // megabytes
		MaxBackups: 50,
		MaxAge:     90, //days
	}))
	envConf := flag.String("env", "config-ci.json", "select a environment config file")
	flag.Parse()
	log.Println("config file ==", *envConf)
	c, err := goconf.New(*envConf)
	if err != nil {
		log.Fatalln("读配置文件出错", err)
	}

	//填充配置文件
	c.Get("/config/LISTEN_PORT", &listenPort)
	c.Get("/config/DB_URL", &dbUrl)
	c.Get("/config/CCDB_URL", &ccdbUrl)

	c.Get("/config/OPEN_CONNS", &openConns)
	c.Get("/config/IDLE_CONNS", &idleConns)

}

func main() {
	dbcomm.InitDB(dbUrl, ccdbUrl, idleConns, openConns)

	go_WebServer()
}
