package pwd

import (
	"encoding/json"
	"fmt"
	"jcd/control/common"
	"jcd/service/dbcomm"

	//	"jcd/service/jcity"
	"jcd/service/jprovince"
	"net/http"
	"strconv"
)

/*
	区域查询请求
*/
type RegionReq struct {
	RegionType int    `json:"region_type"`
	ParentNo   string `json:"parent_no"`
}

/*
	区域查询应答
*/
type RegionResp struct {
	ErrCode string   `json:"err_code"`
	ErrMsg  string   `json:"err_msg"`
	List    []Region `json:"list"`
}

type Region struct {
	RegionNo   string `json:"region_no"`
	RegionName string `json:"region_name"`
}

/*
	说明：得到区域列表
	入参：
	出参：参数1：Token
		 参数1：Error
*/

func GetRegionList(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("GetRegionList")
	userId, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	fmt.Println(uId)
	var regionReq RegionReq
	var regionResp RegionResp

	err := json.NewDecoder(req.Body).Decode(&regionReq)
	if err != nil {
		regionResp.ErrCode = common.ERR_CODE_JSONERR
		regionResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(regionResp, w, req)
		return
	}
	defer req.Body.Close()

	golist := make([]Region, 0)
	switch regionReq.RegionType {
	case common.REGION_PROVINCE:
		var search province.Search
		r := province.New(dbcomm.GetDB(), province.DEBUG)
		l, err := r.GetList(search)
		if err != nil {
			regionResp.ErrCode = common.ERR_CODE_DBERROR
			regionResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR] + err.Error()
			common.Write_Response(regionResp, w, req)
			return
		}
		for _, v := range l {
			var e Region
			e.RegionNo = fmt.Sprintf("%d", v.ProviceId)
			e.RegionName = v.ProviceName
			golist = append(golist, e)
		}
	case common.REGION_CITY:
	case common.REGION_COUNTY:
	case common.REGION_TOWN:
	case common.REGION_VILLAGE:

	}

	regionResp.ErrCode = common.ERR_CODE_SUCCESS
	regionResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS] + "修改密码成功！"
	regionResp.List = golist
	common.Write_Response(regionResp, w, req)
}
