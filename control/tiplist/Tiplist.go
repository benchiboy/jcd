package tiplist

import (
	"encoding/json"
	"jcd/control/common"
	"jcd/service/dbcomm"
	"jcd/service/tiplist"
	"net/http"
)

/*
	查询主页还款及上报请求
*/
type tipListReq struct {
	GetDate  int `json:"get_date"`
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_no"`
}

/*
	查询老赖应答
*/
type tipListResp struct {
	ErrCode string            `json:"err_code"`
	ErrMsg  string            `json:"err_msg"`
	List    []tiplist.Tiplist `json:"list"`
}

/*
	获取老赖指数列表
*/

func GetTipList(w http.ResponseWriter, req *http.Request) {

	var tipReq tipListReq
	var tipResp tipListResp
	err := json.NewDecoder(req.Body).Decode(&tipReq)
	if err != nil {
		tipResp.ErrCode = common.ERR_CODE_JSONERR
		tipResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(tipResp, w, req)
		return
	}
	defer req.Body.Close()

	var search tiplist.Search
	search.PageNo = tipReq.PageNo
	search.PageSize = tipReq.PageSize
	r := tiplist.New(dbcomm.GetDB(), tiplist.DEBUG)
	l, err := r.GetList(search)
	if len(l) == 0 {
		tipResp.ErrCode = common.ERR_CODE_NOTFIND
		tipResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND]
		common.Write_Response(tipResp, w, req)
		return
	}
	tipResp.List = l
	tipResp.ErrCode = common.ERR_CODE_SUCCESS
	tipResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(tipResp, w, req)
}
