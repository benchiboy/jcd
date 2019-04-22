package badloan

import (
	"encoding/json"
	"jcd/control/common"

	"jcd/service/badloan"
	"jcd/service/dbcomm"
	"net/http"
)

/*
	查询老赖请求
*/
type BadLoanListReq struct {
	UserName string `json:"user_name"`
	Phone    string `json:"phone"`
	IdNo     string `json:"id_no"`
	PageNo   int    `json:"page_no"`
	PageSize int    `json:"page_no"`
}

/*
	查询老赖应答
*/
type BadLoanListResp struct {
	ErrCode string            `json:"err_code"`
	ErrMsg  string            `json:"err_msg"`
	Total   int               `json:"total"`
	BadLoan []badloan.Badloan `json:"list"`
}

/*
	获取老赖指数列表
*/

func BadPLoanList(w http.ResponseWriter, req *http.Request) {
	_, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	var listReq BadLoanListReq
	var listResp BadLoanListResp
	err := json.NewDecoder(req.Body).Decode(&listReq)
	if err != nil {
		listResp.ErrCode = common.ERR_CODE_JSONERR
		listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(listResp, w, req)
		return
	}
	defer req.Body.Close()

	var search badloan.Search
	search.PageNo = listReq.PageNo
	search.PageSize = listReq.PageSize
	r := badloan.New(dbcomm.GetDB(), badloan.DEBUG)
	l, err := r.GetList(search)
	total, err := r.GetTotal(search)

	listResp.ErrCode = common.ERR_CODE_SUCCESS
	listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	listResp.BadLoan = l
	listResp.Total = total
	common.Write_Response(listResp, w, req)
}
