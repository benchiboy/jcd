package oauth

import (
	"jcd/control/common"
	"net/http"
)

func Htmllogin(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("Htmllogin")

	common.PrintTail("Htmllogin")
	return
}
