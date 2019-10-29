package payutil

import (
	"testing"
)

//func TestWxOrderQuery(t *testing.T) {
//	WxOrderQuery("1560838209853256000")
//}

//func TestWxOrderQuery(t *testing.T) {
//	WxOrderQuery("1560838209853256001")
//}

//func CPCNPayOrder(bankId string, crfUid string, fullName string, idNo string, cardNo string, mctTrxnNo string, phone string, totalFee int) (string, string, error) {

func TestUnionPayOrder(t *testing.T) {
	CPCNPayOrder("CCB", "5cf18705354168ee6f56feb5294b3164", "杨敬勃", "220202199407164214", "6222800839280003473", "1111A0000000001", "13843276375", 100)
}

func TestUnionPayQueryOrder(t *testing.T) {
	CPCNPayQueryOrder("1111A0000000001")
}
