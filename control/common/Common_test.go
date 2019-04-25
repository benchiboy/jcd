package common

import (
	"testing"
)

func TestGetQrCode(t *testing.T) {
	code, _ := CreateQrCode("23232323233", "weixin://wxpay/bizpayurl?pr=Ab0vBRA")
	t.Log(code)
}
