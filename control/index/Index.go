package index

import (
	"encoding/json"
	"jcd/control/common"

	//	"fmt"
	"jcd/service/dbcomm"
	"jcd/service/index"
	"net/http"
)

/*
	查询老赖请求
*/
type BadIndexListReq struct {
	IndexDate string `json:"index_date"`
	Subject   string `json:"subject"`
	PageNo    int    `json:"page_no"`
	PageSize  int    `json:"page_no"`
}

/*
	查询老赖应答
*/
type BadIndexListResp struct {
	ErrCode string      `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Legend  []string    `json:"legend"`
	XAxis   []XAxisItem `json:"xAxis"`
	YAxis   []YAxisItem `json:"yAxis"`
	Series  []Series    `json:"series"`
}

type XAxisItem struct {
	Type string   `json:"type"`
	Data []string `json:"data"`
}

type YAxisItem struct {
	Type string `json:"type"`
}

type Series struct {
	Name string    `json:"name"`
	Type string    `json:"type"`
	Data []float64 `json:"data"`
}

/*
	获取老赖指数列表
*/

func BadPLoanList(w http.ResponseWriter, req *http.Request) {
	_, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	var listReq BadIndexListReq
	var listResp BadIndexListResp
	err := json.NewDecoder(req.Body).Decode(&listReq)
	if err != nil {
		listResp.ErrCode = common.ERR_CODE_JSONERR
		listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(listResp, w, req)
		return
	}
	defer req.Body.Close()

	var search index.Search
	search.PageNo = listReq.PageNo
	search.PageSize = listReq.PageSize
	r := index.New(dbcomm.GetDB(), index.DEBUG)
	l, err := r.GetList(search)

	lengMap := make(map[string]string, 0)

	legendSlice := make([]string, 0)
	xSlice := make([]XAxisItem, 0)
	ySlice := make([]YAxisItem, 0)
	sSeriesSlice := make([]Series, 0)

	var yi YAxisItem
	yi.Type = "value"
	ySlice = append(ySlice, yi)

	var xi XAxisItem
	xi.Type = "category"

	var si Series
	for _, v := range l {
		lengMap[v.IndexCategory] = v.IndexCategory
		xi.Data = append(xi.Data, v.IndexName)
		si.Data = append(si.Data, v.IndexValue)
	}
	for _, v := range lengMap {
		legendSlice = append(legendSlice, v)
	}
	xSlice = append(xSlice, xi)
	for _, v := range lengMap {
		si.Type = "bar"
		si.Name = v
	}
	sSeriesSlice = append(sSeriesSlice, si)
	listResp.ErrCode = common.ERR_CODE_SUCCESS
	listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	listResp.Legend = legendSlice
	listResp.XAxis = xSlice
	listResp.YAxis = ySlice
	listResp.Series = sSeriesSlice
	common.Write_Response(listResp, w, req)
}
