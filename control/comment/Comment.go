package comment

import (
	"encoding/json"
	"jcd/control/common"

	"fmt"
	"jcd/service/comment"
	"jcd/service/dbcomm"
	"net/http"

	"strconv"
	"time"
)

/*
	获取评论列表
*/
type CommentListReq struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_no"`
}

/*
	获取评论列表返回
*/
type CommentListResp struct {
	ErrCode  string            `json:"err_code"`
	ErrMsg   string            `json:"err_msg"`
	Total    int               `json:"total"`
	CommList []comment.Comment `json:"list"`
}

/*
	喜欢某个评论
*/
type CommentLikeReq struct {
	Id int64 `json:"id"`
}

/*
	喜欢某个评论结果
*/
type CommentLikeResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Likes   int64  `json:"likes"`
}

/*
	举报某个评论
*/
type CommentKillReq struct {
	Id int64 `json:"id"`
}

/*
	举报某个评论结果
*/
type CommentKillResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Kills   int64  `json:"kills"`
}

/*
	投递一评论
*/
type CommentPostReq struct {
	Title   string `json:"title"`
	Context string `json:"context"`
}

/*
	投递一评论结果
*/
type CommentPostResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

/*
	回复一评论请求
*/
type CommentReplyReq struct {
	CommNo  string `json:"comm_no"`
	Title   string `json:"title"`
	Context string `json:"context"`
}

/*
	回复一评论应答
*/
type CommentReplyResp struct {
	Title   string `json:"title"`
	Context string `json:"context"`
}

/*
	获取用户发送的评论列表
*/

func CommentList(w http.ResponseWriter, req *http.Request) {

	_, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	var listReq CommentListReq
	var listResp CommentListResp
	err := json.NewDecoder(req.Body).Decode(&listReq)
	if err != nil {
		listResp.ErrCode = common.ERR_CODE_JSONERR
		listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(listResp, w, req)
		return
	}
	defer req.Body.Close()

	var search comment.Search
	search.PageNo = listReq.PageNo
	search.PageSize = listReq.PageSize
	r := comment.New(dbcomm.GetDB(), comment.DEBUG)
	l, err := r.GetList(search)
	total, err := r.GetTotal(search)

	listResp.ErrCode = common.ERR_CODE_SUCCESS
	listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	listResp.CommList = l
	listResp.Total = total
	common.Write_Response(listResp, w, req)
}

/*
	发表评论
*/

func CommentPost(w http.ResponseWriter, req *http.Request) {
	userId, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	var postReq CommentPostReq
	var postResp CommentPostResp
	err := json.NewDecoder(req.Body).Decode(&postReq)
	if err != nil {
		postResp.ErrCode = common.ERR_CODE_JSONERR
		postResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(postResp, w, req)
		return
	}
	defer req.Body.Close()

	r := comment.New(dbcomm.GetDB(), comment.DEBUG)
	var e comment.Comment
	e.UserId = uId
	e.CommNo = time.Now().UnixNano()
	e.Likes = common.COMMENT_INIT_VALUE
	e.Title = postReq.Title
	e.Context = postReq.Context

	e.InsertTime = time.Now().Format("2006-01-02 15:04:05")

	if err = r.InsertEntity(e, nil); err != nil {
		postResp.ErrCode = common.ERR_CODE_DBERROR
		postResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(postResp, w, req)
		return
	}
	postResp.ErrCode = common.ERR_CODE_SUCCESS
	postResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(postResp, w, req)
}

/*
	LIKE 某个评论
*/
func CommentLike(w http.ResponseWriter, req *http.Request) {
	userId, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	var likeReq CommentLikeReq
	var likeResp CommentLikeResp
	err := json.NewDecoder(req.Body).Decode(&likeReq)
	if err != nil {
		likeResp.ErrCode = common.ERR_CODE_JSONERR
		likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(likeResp, w, req)
		return
	}
	defer req.Body.Close()

	r := comment.New(dbcomm.GetDB(), comment.DEBUG)
	var search comment.Search
	search.Id = likeReq.Id
	e, err := r.Get(search)
	if err != nil {
		likeResp.ErrCode = common.ERR_CODE_NOTFIND
		likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND]
		common.Write_Response(likeResp, w, req)
		return
	}
	likeMap := map[string]interface{}{common.FIELD_LIKES: e.Likes + 1,
		common.FIELD_UPDATE_TIME: time.Now().Format("2006-01-02 15:04:05"),
		common.FIELD_UPDATE_USER: uId}
	err = r.UpdateMap(fmt.Sprintf("%d", likeReq.Id), likeMap, nil)
	if err != nil {
		likeResp.ErrCode = common.ERR_CODE_DBERROR
		likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(likeResp, w, req)
		return
	}
	likeResp.ErrCode = common.ERR_CODE_SUCCESS
	likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	likeResp.Likes = e.Likes + 1

	common.Write_Response(likeResp, w, req)
}

/*
	kILL某个评论
*/
func CommentKill(w http.ResponseWriter, req *http.Request) {
	userId, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	var killReq CommentKillReq
	var killResp CommentKillResp
	err := json.NewDecoder(req.Body).Decode(&killReq)
	if err != nil {
		killResp.ErrCode = common.ERR_CODE_JSONERR
		killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(killResp, w, req)
		return
	}
	defer req.Body.Close()

	r := comment.New(dbcomm.GetDB(), comment.DEBUG)
	var search comment.Search
	search.Id = killReq.Id
	e, err := r.Get(search)
	if err != nil {
		killResp.ErrCode = common.ERR_CODE_NOTFIND
		killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND]
		common.Write_Response(killResp, w, req)
		return
	}
	killMap := map[string]interface{}{common.FIELD_KILLS: e.Kills + 1,
		common.FIELD_UPDATE_TIME: time.Now().Format("2006-01-02 15:04:05"),
		common.FIELD_UPDATE_USER: uId}
	err = r.UpdateMap(fmt.Sprintf("%d", killReq.Id), killMap, nil)
	if err != nil {
		killResp.ErrCode = common.ERR_CODE_DBERROR
		killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(killResp, w, req)
		return
	}
	killResp.ErrCode = common.ERR_CODE_SUCCESS
	killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	killResp.Kills = e.Kills + 1

	common.Write_Response(killResp, w, req)
}
