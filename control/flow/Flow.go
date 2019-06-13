package flow

import (
	"encoding/json"
	"jcd/control/common"

	"jcd/service/dbcomm"
	"jcd/service/disputes"

	"jcd/service/flow"

	"jcd/service/bizaccount"
	"jcd/service/bizcustomer"

	"fmt"
	//	"jcd/control/payutil"
	"net/http"
	"strconv"
	"time"
)

/*
	查询输入条件
*/
type FindLoanReq struct {
	KeyNo       string `json:"key_no"`
	IdKey       string `json:"id_key"`
	CapchasCode string `json:"captchas_code"`
	PageNo      int    `json:"page_no"`
	PageSize    int    `json:"page_no"`
}

/*
	查询输出条件
*/
type FindLoanResp struct {
	ErrCode  string `json:"err_code"`
	ErrMsg   string `json:"err_msg"`
	Total    int    `json:"total"`
	LoanList []Loan `json:"loan_list"`
}

type Loan struct {
	MctNo        string `json:"mct_no"`
	MainAddr     string `json:"main_addr"`
	MctName      string `json:"mct_name"`
	CertNo       string `json:"cert_no"`
	CustName     string `json:"cust_name"`
	Phone        string `json:"phone"`
	LoanNo       string `json:"loan_no"`
	PayDate      string `json:"pay_date"`
	PayAmt       string `json:"pay_amt"`
	LoanDay      string `json:"loan_day"`
	DueDate      string `json:"due_Date"`
	OverdueDay   string `json:"overdue_day"`
	TotalAmt     string `json:"total_amt"`
	PrincipalAmt string `json:"principal_amt"`
	IntFeeAmt    string `json:"intfee_amt"`
}

/*
   发起还款的请求
*/
type RepayReq struct {
	MctNo     string `json:"mct_no"`
	TradeType string `json:"trade_type"`
	LoanNo    string `json:"loan_no"`
	RepayAmt  string `json:"repay_amt"`
}

/*
   还款的应答
*/
type RepayResp struct {
	QrCode  string `json:"qr_code"`
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
   对逾期待还款意见反馈
*/
type DisputesReq struct {
	MctNo        string `json:"mct_no"`
	LoanNo       string `json:"loan_no"`
	Phone        string `json:"phone"`
	Mail         string `json:"mail"`
	DisputesMemo string `json:"feedback"`
}

/*
   还款的应答
*/
type DisputesResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
	获取用户的交易流水请求
*/
type FlowListReq struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

/*
	获取用户的交易流水列表
*/
type FlowListResp struct {
	ErrCode  string      `json:"err_code"`
	ErrMsg   string      `json:"err_msg"`
	Total    int         `json:"total"`
	FlowList []flow.Flow `json:"list"`
}

/*
	获取用户的交易流水请求
*/
type FeedbackReq struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

/*
	获取用户的交易流水列表
*/
type FeedbackListResp struct {
	ErrCode      string              `json:"err_code"`
	ErrMsg       string              `json:"err_msg"`
	Total        int                 `json:"total"`
	DisputesList []disputes.Disputes `json:"list"`
}

/*
	修改当前登录用户的密码
	当前机构及用户ID从TOKEN获取
*/
func Disputes(w http.ResponseWriter, req *http.Request) {
	//	userId, _, tokenErr := common.CheckToken(w, req)
	//	if tokenErr != nil {
	//		return
	//	}
	var disputesReq DisputesReq
	var disputesResp DisputesResp
	err := json.NewDecoder(req.Body).Decode(&disputesReq)
	if err != nil {
		disputesResp.ErrCode = common.ERR_CODE_JSONERR
		disputesResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(disputesResp, w, req)
		return
	}
	defer req.Body.Close()

	var search disputes.Search
	search.MctNo = disputesReq.MctNo
	search.MctLoanNo = disputesReq.LoanNo
	r := disputes.New(dbcomm.GetDB(), disputes.DEBUG)
	if _, err := r.Get(search); err == nil {
		disputesResp.ErrCode = common.ERR_CODE_EXISTED
		disputesResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED]
		common.Write_Response(disputesResp, w, req)
		return
	} else {
		var e disputes.Disputes
		//		uid, _ := strconv.ParseInt(userId, 10, 64)
		//		e.UserId = uid
		e.MctNo = disputesReq.MctNo
		e.InsertTime = time.Now().Format("2006-01-02 15:04:05")
		e.MctLoanNo = disputesReq.LoanNo
		e.DisputeNo = time.Now().UnixNano()
		e.Mail = disputesReq.Mail
		e.Phone = disputesReq.Phone
		e.Status = common.STATUS_INIT
		e.DisputeMemo = disputesReq.DisputesMemo
		r.InsertEntity(e, nil)
	}
	disputesResp.ErrCode = common.ERR_CODE_SUCCESS
	disputesResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(disputesResp, w, req)

}

/*
	查询用户的待还记录
*/
func FindLoans(w http.ResponseWriter, req *http.Request) {
	var searchReq FindLoanReq
	var searchResp FindLoanResp
	err := json.NewDecoder(req.Body).Decode(&searchReq)
	if err != nil {
		searchResp.ErrCode = common.ERR_CODE_JSONERR
		searchResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(searchResp, w, req)
		return
	}
	defer req.Body.Close()

	//验证图形验证码
	fmt.Println("code==>", searchReq.CapchasCode, "idkey===>", searchReq.IdKey)
	if !common.CheckCaptchaCode(searchReq.IdKey, searchReq.CapchasCode) {
		searchResp.ErrCode = common.ERR_CODE_VERIFY
		searchResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_VERIFY]
		common.Write_Response(searchResp, w, req)
		return
	}

	var search bizcustomer.Search
	search.CertNo = searchReq.KeyNo
	c := bizcustomer.New(dbcomm.GetCCDB(), bizcustomer.DEBUG)
	l, err := c.GetList(search)
	for _, v := range l {
		var acctSearch bizaccount.Search
		acctSearch.CrfUid = v.CrfUid
		a := bizaccount.New(dbcomm.GetCCDB(), bizaccount.DEBUG)
		ll, _ := a.GetList(acctSearch)
		for _, vv := range ll {
			var e Loan
			e.CertNo = v.CertNo
			e.MctNo = "CRF01"
			e.MainAddr = v.MainContactAddr
			e.CustName = v.CustName
			e.Phone = v.MCardMobile
			e.LoanNo = vv.ContractId
			e.PayAmt = vv.LoanCapital
			e.DueDate = vv.DueDate
			e.IntFeeAmt = vv.IntFeeAmt
			e.OverdueDay = vv.OverdueDays
			e.TotalAmt = vv.Total
			e.PrincipalAmt = vv.Pricipal
			e.PayDate = vv.PayDate
			searchResp.LoanList = append(searchResp.LoanList, e)
		}
	}
	searchResp.ErrCode = common.ERR_CODE_SUCCESS
	searchResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(searchResp, w, req)

}

/*
	用户发起还款
*/
func RepayOrder(w http.ResponseWriter, req *http.Request) {
	userId, _, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	var repayReq RepayReq
	var repayResp RepayResp
	err := json.NewDecoder(req.Body).Decode(&repayReq)
	if err != nil {
		repayResp.ErrCode = common.ERR_CODE_JSONERR
		repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(repayResp, w, req)
		return
	}
	defer req.Body.Close()
	var qrCode string

	//	qrCode, err = common.CreateQrCode("23232232323232", "weixin://wxpay/bizpayurl?pr=B5GobHj")
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		repayResp.ErrCode = common.ERR_CODE_PAYERR
	//		repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_PAYERR]
	//		common.Write_Response(repayResp, w, req)
	//		return
	//	}
	//	repayResp.QrCode = qrCode
	//	repayResp.ErrCode = common.ERR_CODE_SUCCESS
	//	repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	//	common.Write_Response(repayResp, w, req)
	//	return

	var search flow.Search
	search.MctNo = repayReq.MctNo
	search.MctTrxnNo = repayReq.LoanNo
	r := flow.New(dbcomm.GetDB(), flow.DEBUG)
	//	if _, err := r.Get(search); err == nil {
	//		repayResp.ErrCode = common.ERR_CODE_EXISTED
	//		repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED]
	//		common.Write_Response(repayResp, w, req)
	//		return

	//	} else {
	var e flow.Flow
	e.MctNo = repayReq.MctNo
	e.TrxnAmt = 1
	uid, _ := strconv.ParseInt(userId, 10, 64)
	e.UserId = uid
	e.MctTrxnNo = repayReq.LoanNo
	e.InsertTime = time.Now().Format("2006-01-02 15:04:05")
	e.TrxnDate = time.Now().Format("2006-01-02 15:04:05")
	e.TrxnNo = time.Now().UnixNano()
	e.ProcStatus = common.STATUS_INIT
	if err := r.InsertEntity(e, nil); err != nil {
		repayResp.ErrCode = common.ERR_CODE_DBERROR
		repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(repayResp, w, req)
		return
	}
	//prePayId, codeUrl, err := payutil.UnionPayOrder(fmt.Sprintf("%d", e.TrxnNo), e.TrxnAmt)
	prePayId := fmt.Sprintf("%d", e.TrxnNo)
	codeUrl := "weixin://wxpay/bizpayurl?pr=B5GobHj"
	if err != nil {
		flowMap := map[string]interface{}{common.FIELD_PROC_STATUS: common.STATUS_FAIL,
			common.FIELD_PROC_MSG: err.Error()}
		err = r.UpdateMap(fmt.Sprintf("%d", e.TrxnNo), flowMap, nil)
		if err != nil {
			repayResp.ErrCode = common.ERR_CODE_DBERROR
			repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		} else {
			repayResp.ErrCode = common.ERR_CODE_PAYERR
			repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_PAYERR]
		}
		common.Write_Response(repayResp, w, req)
		return

	} else {
		flowMap := map[string]interface{}{common.FIELD_PROC_STATUS: common.STATUS_DOING,
			common.FIELD_PREPAY_ID: prePayId,
			common.FIELD_CODE_URL:  codeUrl}
		err = r.UpdateMap(fmt.Sprintf("%d", e.TrxnNo), flowMap, nil)
		if err != nil {
			repayResp.ErrCode = common.ERR_CODE_PAYERR
			repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_PAYERR]
			common.Write_Response(repayResp, w, req)
			return
		}
		qrCode, err = common.CreateQrCode(prePayId, codeUrl)
		if err != nil {
			fmt.Println(err.Error())
			repayResp.ErrCode = common.ERR_CODE_PAYERR
			repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_PAYERR]
			common.Write_Response(repayResp, w, req)
			return
		}
	}

	repayResp.QrCode = qrCode
	repayResp.ErrCode = common.ERR_CODE_SUCCESS
	repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(repayResp, w, req)
}

/*
	获取用户的还款交易流水
*/
func MyFlowList(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("MyFlowList")
	userId, _, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	var listReq FlowListReq
	var listResp FlowListResp
	err := json.NewDecoder(req.Body).Decode(&listReq)
	if err != nil {
		listResp.ErrCode = common.ERR_CODE_JSONERR
		listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(listResp, w, req)
		return
	}
	defer req.Body.Close()
	var search flow.Search
	search.PageNo = listReq.PageNo
	search.PageSize = listReq.PageSize
	search.UserId = uId
	r := flow.New(dbcomm.GetDB(), flow.DEBUG)
	l, err := r.GetList(search)
	total, err := r.GetTotal(search)
	listResp.ErrCode = common.ERR_CODE_SUCCESS
	listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	listResp.FlowList = l
	listResp.Total = total
	common.Write_Response(listResp, w, req)
	common.PrintTail("MyFlowList")
}

/*
	获取用户的还款交易流水
*/
func MyFeedbackList(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("MyFeedbackList")
	userId, _, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	var listReq FeedbackReq
	var listResp FeedbackListResp
	err := json.NewDecoder(req.Body).Decode(&listReq)
	if err != nil {
		listResp.ErrCode = common.ERR_CODE_JSONERR
		listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(listResp, w, req)
		return
	}
	defer req.Body.Close()
	var search disputes.Search
	search.PageNo = listReq.PageNo
	search.PageSize = listReq.PageSize
	search.UserId = uId
	r := disputes.New(dbcomm.GetDB(), disputes.DEBUG)
	l, err := r.GetList(search)
	total, err := r.GetTotal(search)
	listResp.ErrCode = common.ERR_CODE_SUCCESS
	listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	listResp.DisputesList = l
	listResp.Total = total
	common.Write_Response(listResp, w, req)
	common.PrintTail("MyFeedbackList")
}
