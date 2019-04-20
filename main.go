// rcs_contract_mgr project main.go
package main

import (
	"flag"

	"io"

	"jcd/control/account"
	"jcd/control/flow"
	"jcd/control/login"
	"jcd/control/pwd"
	"jcd/control/smscode"
	"log"
	"net/http"
	"os"

	"jcd/service/dbcomm"

	goconf "github.com/pantsing/goconf"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	http_srv   *http.Server
	dbUrl      string
	listenPort int
	idleConns  int
	openConns  int
)

func go_WebServer() {
	log.Println("Listen Service start...")

	http.HandleFunc("/wxLogin", login.WxLogin)

	http.HandleFunc("/jc/api/disputes", flow.Disputes)
	http.HandleFunc("/jc/api/signin", login.SignIn)
	http.HandleFunc("/jc/api/signout", login.SignOut)
	http.HandleFunc("/jc/api/signup", account.SignUp)
	http.HandleFunc("/jc/api/chgpwd", pwd.ChangePwd)
	http.HandleFunc("/jc/api/callback", account.GetAccount)
	http.HandleFunc("/jc/api/findloan", flow.FindLoans)
	http.HandleFunc("/jc/api/repay", flow.RepayOrder)
	http.HandleFunc("/jc/api/getsmscode", smscode.GetSmsCode)
	http.HandleFunc("/jc/api/checksmscode", smscode.CheckSmsCode)
	http.HandleFunc("/jc/api/getcaptcha", smscode.GetCaptchas)

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
	c.Get("/config/OPEN_CONNS", &openConns)
	c.Get("/config/IDLE_CONNS", &idleConns)

}

func main() {
	dbcomm.InitDB(dbUrl, idleConns, openConns)
	go_WebServer()
}
