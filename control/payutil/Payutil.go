package payutil

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"jcd/control/common"
	"jcd/service/dbcomm"
	"jcd/service/flow"

	"net/http"
	"sort"
	"strings"
	"time"
)

//首先定义一个UnifyOrderReq用于填入我们要传入的参数。
type UnifyOrderReq struct {
	Appid            string `xml:"appid"`
	Body             string `xml:"body"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Trade_type       string `xml:"trade_type"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        int    `xml:"total_fee"`
	Out_trade_no     string `xml:"out_trade_no"`
	Sign             string `xml:"sign"`
}

type UnifyOrderResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Nonce       string `xml:"nonce_str"`
	Result_code string `xml:"result_code"`
	Trade_type  string `xml:"trade_type"`
	Prepay_id   string `xml:"prepay_id"`
	Code_url    string `xml:"code_url"`
}

type WXPayNotifyReq struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce          string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Openid         string `xml:"openid"`
	Is_subscribe   string `xml:"is_subscribe"`
	Trade_type     string `xml:"trade_type"`
	Bank_type      string `xml:"bank_type"`
	Total_fee      int    `xml:"total_fee"`
	Fee_type       string `xml:"fee_type"`
	Cash_fee       int    `xml:"cash_fee"`
	Cash_fee_Type  string `xml:"cash_fee_type"`
	Transaction_id string `xml:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	Time_end       string `xml:"time_end"`
}
type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

//微信支付计算签名的函数
func wxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

/*
	得到支付订单的状态
*/
func GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("GetOrderStatus")

	common.PrintTail("GetOrderStatus")
}

func UnionPayOrder(mctTrxnNo string, totalFee int) (string, string, error) {
	var orderReq UnifyOrderReq
	orderReq.Appid = common.APP_ID
	orderReq.Body = common.PRODUCT_NAME
	orderReq.Mch_id = common.MCT_ID
	orderReq.Nonce_str = fmt.Sprintf("%d", time.Now().Unix())
	orderReq.Notify_url = common.WX_PAY_CALLBACK_URL
	orderReq.Trade_type = common.TRADE_TYPE_NATIVE
	orderReq.Spbill_create_ip = common.SERVER_IP
	orderReq.Total_fee = totalFee
	orderReq.Out_trade_no = mctTrxnNo

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = orderReq.Appid
	m["body"] = orderReq.Body
	m["mch_id"] = orderReq.Mch_id
	m["notify_url"] = orderReq.Notify_url
	m["trade_type"] = orderReq.Trade_type
	m["spbill_create_ip"] = orderReq.Spbill_create_ip
	m["total_fee"] = orderReq.Total_fee
	m["out_trade_no"] = orderReq.Out_trade_no
	m["nonce_str"] = orderReq.Nonce_str
	orderReq.Sign = wxpayCalcSign(m, common.MCT_KEY) //这个是计算wxpay签名的函数上面已贴出
	bytes_req, err := xml.Marshal(orderReq)
	if err != nil {
		fmt.Println("以xml形式编码发送错误, 原因:", err)
		return common.EMPTY_STRING, common.EMPTY_STRING, err
	}
	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "UnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)
	fmt.Println(string(bytes_req))
	//发送unified order请求.
	req, err := http.NewRequest("POST", common.WX_PAY_URL, bytes.NewReader(bytes_req))
	if err != nil {
		fmt.Println("New Http Request发生错误，原因:", err)
		return common.EMPTY_STRING, common.EMPTY_STRING, err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	c := http.Client{}
	resp, _err := c.Do(req)
	if _err != nil {
		fmt.Println("请求微信支付统一下单接口发送错误, 原因:", _err)
		return common.EMPTY_STRING, common.EMPTY_STRING, err
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return common.EMPTY_STRING, common.EMPTY_STRING, err
	}
	fmt.Println("Body=================>", string(body))
	var uniResp UnifyOrderResp
	err = xml.Unmarshal(body, &uniResp)
	fmt.Println("=====================>")
	fmt.Println("=====================>", uniResp.Code_url)
	return uniResp.Prepay_id, uniResp.Code_url, nil
}

//微信支付签名验证函数
func wxpayVerifySign(needVerifyM map[string]interface{}, sign string) bool {
	signCalc := wxpayCalcSign(needVerifyM, common.MCT_KEY)
	fmt.Printf("计算出来的sign: %v", signCalc)
	fmt.Printf("微信异步通知sign: %v", sign)
	if sign == signCalc {
		fmt.Println("签名校验通过!")
		return true
	}
	fmt.Println("签名校验失败!")
	return false
}

//具体的微信支付回调函数的范例
func WxpayCallback(w http.ResponseWriter, r *http.Request) {
	common.PrintHead("WxpayCallback")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取http body失败，原因!", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	fmt.Println("微信支付异步通知，HTTP Body:", string(body))
	var mr WXPayNotifyReq
	err = xml.Unmarshal(body, &mr)
	if err != nil {
		fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)
	reqMap["return_code"] = mr.Return_code
	reqMap["return_msg"] = mr.Return_msg
	reqMap["appid"] = mr.Appid
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce
	reqMap["result_code"] = mr.Result_code
	reqMap["openid"] = mr.Openid
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["trade_type"] = mr.Trade_type
	reqMap["bank_type"] = mr.Bank_type
	reqMap["total_fee"] = mr.Total_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["cash_fee_type"] = mr.Cash_fee_Type
	reqMap["transaction_id"] = mr.Transaction_id
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["attach"] = mr.Attach
	reqMap["time_end"] = mr.Time_end
	var resp WXPayNotifyResp
	if wxpayVerifySign(reqMap, mr.Sign) {
		fmt.Println("微信支付成功....")
		resp.Return_code = "SUCCESS"
		resp.Return_msg = "OK"
		r := flow.New(dbcomm.GetDB(), flow.DEBUG)
		resetMap := map[string]interface{}{common.FIELD_PROC_STATUS: common.STATUS_SUCC,
			common.FIELD_PROC_MSG: common.SUCC_MSG}
		r.UpdateMap(mr.Out_trade_no, resetMap, nil)

	} else {
		resp.Return_code = "FAIL"
		resp.Return_msg = "failed to verify sign, please retry!"
	}
	//结果返回，微信要求如果成功需要返回return_code "SUCCESS"
	bytes, _err := xml.Marshal(resp)
	strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
	if _err != nil {
		fmt.Println("xml编码失败，原因：", _err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	fmt.Fprint(w.(http.ResponseWriter), strResp)
}
