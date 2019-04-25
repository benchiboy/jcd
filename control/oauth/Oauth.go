package oauth

import (
	"fmt"
	"io/ioutil"
	"jcd/control/common"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

var endpotin = oauth2.Endpoint{
	AuthURL:  "https://api.weibo.com/OAuth2/authorize",
	TokenURL: "https://api.weibo.com/OAuth2/access_token",
}

var oAuthConfig = &oauth2.Config{
	ClientID:     "2302557195",
	ClientSecret: "3e096effd83cd4d0553b5098159eced8",
	RedirectURL:  "http://132.232.11.85:8087/jc/api/weibocallback",
	Scopes:       []string{"https://api.weibo.com/OAuth2/access_token"},
	Endpoint:     endpotin,
}

const oauthStateString = "random"

func WeiboCallback(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("WeiboCallback")
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	fmt.Println("WeiboCallback======>", code)
	token, err := oAuthConfig.Exchange(oauth2.NoContext, code)
	fmt.Println(token)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	response, err := http.PostForm("https://api.weibo.com/Oauth2/get_token_info", url.Values{"access_token": {token.AccessToken}})
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)

	common.PrintTail("WeiboCallback")

	return
}

func WeiboCancelCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("......proc user reg......")
	t1 := time.Now()
	log.Println(t1)
	return
}

func WeiboLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("......proc user reg......")
	t1 := time.Now()
	url := oAuthConfig.AuthCodeURL(oauthStateString)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	log.Println("Create user elapsed time:", time.Since(t1))
	return
}

func WeiboLoginIndex(w http.ResponseWriter, r *http.Request) {

	htmlIndex := `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
</body></html>`
	fmt.Fprintf(w, htmlIndex)

}

const htmlIndex = `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
</body></html>
`
