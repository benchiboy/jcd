package flow

import (
	"encoding/json"
	"jcd/control/common"

	"jcd/service/dbcomm"
	"jcd/service/disputes"
	"jcd/service/flow"

	"fmt"
	"jcd/control/payutil"
	"net/http"
	"strconv"
	"time"
)

/*
	查询输入条件
*/
type FindLoanReq struct {
	KeyNo    string `json:"key_no"`
	PageNo   int    `json:"page_no"`
	PageSize int    `json:"page_no"`
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
	MctName      string `json:"mct_name"`
	LoanNo       string `json:"loan_no"`
	PayDate      string `json:"pay_date"`
	PayAmt       string `json:"pay_amt"`
	LoanDay      string `json:"loan_day"`
	DueDate      string `json:"due_Date"`
	OverdueDay   string `json:"overdue_day"`
	TotalBalAmt  string `json:"total_balamt"`
	PrincipalAmt string `json:"pricipal_balamt"`
	IntFeeAmt    string `json:"intfee_balamt"`
}

/*
   发起还款的请求
*/
type RepayReq struct {
	MctNo     string `json:"mct_no"`
	TradeType string `json:"trade_type"`
	LoanNo    string `json:"loan_no"`
	RepayAmt  int    `json:"repay_amt"`
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
   发起还款的请求
   type :1 还款，2: 还款异议
*/
type DisputesReq struct {
	MctNo        string `json:"mct_no"`
	LoanNo       string `json:"loan_no"`
	RepayAmt     string `json:"repay_amt"`
	DisputesMemo string `json:"memo"`
}

/*
   还款的应答
*/
type DisputesResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
	修改当前登录用户的密码
	当前机构及用户ID从TOKEN获取
*/
func Disputes(w http.ResponseWriter, req *http.Request) {
	userId, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
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
	search.MctTrxnNo = disputesReq.LoanNo
	r := disputes.New(dbcomm.GetDB(), disputes.DEBUG)
	if _, err := r.Get(search); err == nil {
		disputesResp.ErrCode = common.ERR_CODE_EXISTED
		disputesResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED]
		common.Write_Response(disputesResp, w, req)
		return

	} else {
		var e disputes.Disputes

		uid, _ := strconv.ParseInt(userId, 10, 64)
		e.UserId = uid
		e.MctNo = disputesReq.MctNo
		f, _ := strconv.ParseFloat(disputesReq.RepayAmt, 32)
		e.DisputeAmt = f
		e.InsertTime = time.Now().Format("2006-01-02 15:04:05")
		e.DisputeDate = time.Now().Format("2006-01-02 15:04:05")
		e.MctTrxnNo = disputesReq.LoanNo
		e.DisputeNo = time.Now().UnixNano()
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
	_, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
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

	var e Loan
	e.LoanNo = "1809999199191"
	e.LoanDay = "12"
	e.MctNo = "100000"
	e.MctName = "信而富-QQ现金贷"
	e.PayAmt = "1000.00"
	e.PayDate = "2018-12-12"
	e.OverdueDay = "100"
	e.PrincipalAmt = "2000.00"
	e.IntFeeAmt = "200.00"
	e.TotalBalAmt = "20200.00"

	searchResp.ErrCode = common.ERR_CODE_SUCCESS
	searchResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	searchResp.LoanList = append(searchResp.LoanList, e)
	searchResp.LoanList = append(searchResp.LoanList, e)
	searchResp.Total = len(searchResp.LoanList)
	common.Write_Response(searchResp, w, req)

}

/*
	用户发起还款
*/
func RepayOrder(w http.ResponseWriter, req *http.Request) {
	userId, _, tokenErr := common.CheckToken(w, req)
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
	var search flow.Search
	search.MctNo = repayReq.MctNo
	search.MctTrxnNo = repayReq.LoanNo
	r := flow.New(dbcomm.GetDB(), flow.DEBUG)
	if _, err := r.Get(search); err == nil {
		repayResp.ErrCode = common.ERR_CODE_EXISTED
		repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED]
		common.Write_Response(repayResp, w, req)
		return

	} else {
		var e flow.Flow
		e.MctNo = repayReq.MctNo
		e.TrxnAmt = repayReq.RepayAmt
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
		prePayId, codeUrl, err := payutil.UnionPayOrder(fmt.Sprintf("%d", e.TrxnNo), e.TrxnAmt)
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
	}
	repayResp.QrCode = qrCode
	repayResp.ErrCode = common.ERR_CODE_SUCCESS
	repayResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(repayResp, w, req)

}
