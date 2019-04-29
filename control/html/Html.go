package oauth

import (
	"fmt"
	"jcd/control/common"
	"net/http"
)

func HtmlHome(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("Htmllogin")

	htmlIndex := `<html><body>
<a href="/jc/api/weibologin">Welcome to login WEIBO</a>
</body></html>`
	fmt.Fprintf(w, htmlIndex)

	common.PrintTail("Htmllogin")
	return
}
