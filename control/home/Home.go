package home

import (
	"fmt"
	"jcd/control/common"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("Home")

	htmlIndex := `<html><body>
	<a href="/jc/api/weibologin">Welcome to login WEIBO</a>
		</body></html>`
	fmt.Fprintf(w, htmlIndex)

	common.PrintTail("Home")
	return
}
