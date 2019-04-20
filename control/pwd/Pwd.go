package pwd

import (
	"encoding/json"
	"jcd/control/common"

	"fmt"
	"jcd/service/account"
	"jcd/service/dbcomm"
	"net/http"
)

/*
	修改密码
*/
type ChangePwdReq struct {
	OldPwd string `json:"old_password"`
	NewPwd string `json:"new_password"`
}

/*
 */
type ChangePwdResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
	重置密码
*/
type ResetPwdReq struct {
	FlowBatchId string `json:"flow_batch_id"`
	PassWord    string `json:"password"`
}

/*
 */
type ResetPwdResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
	修改当前登录用户的密码
	当前机构及用户ID从TOKEN获取
*/
func ChangePwd(w http.ResponseWriter, req *http.Request) {
	userId, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	var changeReq ChangePwdReq
	var changeResp ChangePwdResp
	err := json.NewDecoder(req.Body).Decode(&changeReq)
	if err != nil {
		changeResp.ErrCode = common.ERR_CODE_JSONERR
		changeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(changeResp, w, req)
		return
	}

	defer req.Body.Close()

	var search account.Search
	search.LoginName = userId
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	p, err := r.Get(search)
	if err != nil {
		changeResp.ErrCode = common.ERR_CODE_DBERROR
		changeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR] + "修改密码失败" + err.Error()
		common.Write_Response(changeResp, w, req)
		return
	}

	if p.LoginPass != changeReq.OldPwd {
		changeResp.ErrCode = common.ERR_CODE_NOMATCH
		changeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOMATCH] + "旧密码不正确，无法修改密码"
		common.Write_Response(changeResp, w, req)
		return
	}

	if changeReq.NewPwd == changeReq.OldPwd {
		changeResp.ErrCode = common.ERR_CODE_NOMATCH
		changeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOMATCH] + "旧密码和新密码一致，无法修改密码"
		common.Write_Response(changeResp, w, req)
		return
	}
	var p2 account.Account
	p2.LoginPass = changeReq.NewPwd
	err = r.UpdataEntity(fmt.Sprintf("%d", p.Id), p2, nil)
	if err != nil {
		changeResp.ErrCode = common.ERR_CODE_DBERROR
		changeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR] + err.Error()
		common.Write_Response(changeResp, w, req)
		return
	}
	changeResp.ErrCode = common.ERR_CODE_SUCCESS
	changeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS] + "修改密码成功！"
	common.Write_Response(changeResp, w, req)
}

/*
	重置某个用户的密码
	注意：需要再系统记录用户的退出时间，暂时先MOCK
*/
func ResetPwd(w http.ResponseWriter, req *http.Request) {
	userId, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	var resetReq ResetPwdReq
	var resetResp ResetPwdResp
	err := json.NewDecoder(req.Body).Decode(&resetReq)
	if err != nil {
		resetResp.ErrCode = common.ERR_CODE_JSONERR
		resetResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + ":" + err.Error()
		common.Write_Response(resetResp, w, req)
		return
	}
	defer req.Body.Close()

	var search account.Search
	search.LoginName = userId
	r := account.New(dbcomm.GetDB(), account.DEBUG)
	p, err := r.Get(search)
	if err != nil {
		resetResp.ErrCode = common.ERR_CODE_NOTFIND
		resetResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND] + "被重置密码用户不存在！"
		common.Write_Response(resetResp, w, req)
		return
	}
	resetMap := map[string]interface{}{common.FIELD_LOGIN_PASS: common.DEFAULT_PWD, common.FIELD_ERRORS: 0}
	err = r.UpdateMap(fmt.Sprintf("%d", p.Id), resetMap, nil)
	if err != nil {
		resetResp.ErrCode = common.ERR_CODE_DBERROR
		resetResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(resetResp, w, req)
		return
	}
	resetResp.ErrCode = common.ERR_CODE_SUCCESS
	resetResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS] + "重置密码成功！"
	common.Write_Response(resetResp, w, req)
}
