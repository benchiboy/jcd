package index

import (
	"encoding/json"
	"jcd/control/common"
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
	if len(l) == 0 {
		listResp.ErrCode = common.ERR_CODE_NOTFIND
		listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND]
		common.Write_Response(listResp, w, req)
		return
	}

	legendSlice := make([]string, 0)
	xSlice := make([]XAxisItem, 0)
	var bTag bool
	var category string
	category = l[0].IndexCategory
	legendSlice = append(legendSlice, l[0].IndexCategory)

	var xi XAxisItem
	for _, v := range l {
		if v.IndexCategory != category {
			legendSlice = append(legendSlice, v.IndexCategory)
			category = v.IndexCategory
			bTag = true
		} else {
			if !bTag {
				xi.Type = "category"
				xi.Data = append(xi.Data, v.IndexName)
			}
		}
	}

	xSlice = append(xSlice, xi)
	var yi YAxisItem
	ySlice := make([]YAxisItem, 0)
	yi.Type = "value"
	ySlice = append(ySlice, yi)

	sSeriesSlice := make([]Series, 0)
	category = common.EMPTY_STRING
	for _, v := range legendSlice {
		var si Series
		for _, v1 := range l {
			if v1.IndexCategory == v {
				si.Data = append(si.Data, v1.IndexValue)
			}
		}
		si.Name = v
		si.Type = "bar"
		sSeriesSlice = append(sSeriesSlice, si)
	}
	listResp.ErrCode = common.ERR_CODE_SUCCESS
	listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	listResp.Legend = legendSlice
	listResp.XAxis = xSlice
	listResp.YAxis = ySlice
	listResp.Series = sSeriesSlice
	common.Write_Response(listResp, w, req)
}
