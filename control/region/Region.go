package region

import (
	"encoding/json"
	"fmt"
	"jcd/control/common"
	"jcd/service/dbcomm"
	"jcd/service/jcity"
	"jcd/service/jcounty"
	"jcd/service/jprovince"
	"jcd/service/jtown"
	"jcd/service/jvillage"
	"net/http"
	"strconv"
)

/*
	区域查询请求
*/
type RegionReq struct {
	RegionType int   `json:"region_type"`
	ParentNo   int64 `json:"parent_no"`
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
	RegionNo   int64  `json:"region_no"`
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
	userId, _, _, tokenErr := common.CheckToken(w, req)
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
			e.RegionNo = v.ProviceId
			e.RegionName = v.ProviceName
			golist = append(golist, e)
		}
	case common.REGION_CITY:
		var search city.Search
		search.ProvinceId = regionReq.ParentNo
		r := city.New(dbcomm.GetDB(), city.DEBUG)
		l, err := r.GetList(search)
		if err != nil {
			regionResp.ErrCode = common.ERR_CODE_DBERROR
			regionResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR] + err.Error()
			common.Write_Response(regionResp, w, req)
			return
		}
		for _, v := range l {
			var e Region
			e.RegionNo = v.CityId
			e.RegionName = v.CityName
			golist = append(golist, e)
		}
	case common.REGION_COUNTY:
		var search county.Search
		search.CityId = regionReq.ParentNo
		r := county.New(dbcomm.GetDB(), county.DEBUG)
		l, err := r.GetList(search)
		if err != nil {
			regionResp.ErrCode = common.ERR_CODE_DBERROR
			regionResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR] + err.Error()
			common.Write_Response(regionResp, w, req)
			return
		}
		for _, v := range l {
			var e Region
			e.RegionNo = v.CountyId
			e.RegionName = v.CountyName
			golist = append(golist, e)
		}
	case common.REGION_TOWN:
		var search town.Search
		search.CountyId = regionReq.ParentNo
		r := town.New(dbcomm.GetDB(), town.DEBUG)
		l, err := r.GetList(search)
		if err != nil {
			regionResp.ErrCode = common.ERR_CODE_DBERROR
			regionResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR] + err.Error()
			common.Write_Response(regionResp, w, req)
			return
		}
		for _, v := range l {
			var e Region
			e.RegionNo = v.TownId
			e.RegionName = v.TownName
			golist = append(golist, e)
		}
	case common.REGION_VILLAGE:
		var search village.Search
		search.TownId = regionReq.ParentNo
		r := village.New(dbcomm.GetDB(), village.DEBUG)
		l, err := r.GetList(search)
		if err != nil {
			regionResp.ErrCode = common.ERR_CODE_DBERROR
			regionResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR] + err.Error()
			common.Write_Response(regionResp, w, req)
			return
		}
		for _, v := range l {
			var e Region
			e.RegionNo = v.VillageId
			e.RegionName = v.VillageName
			golist = append(golist, e)
		}

	}
	regionResp.ErrCode = common.ERR_CODE_SUCCESS
	regionResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	regionResp.List = golist
	common.Write_Response(regionResp, w, req)
}
