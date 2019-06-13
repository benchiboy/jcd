package comment

import (
	"encoding/json"
	"jcd/control/common"

	"fmt"
	"jcd/service/comment"
	"jcd/service/commentuser"
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
	PageSize int `json:"page_size"`
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
	CommNo int64 `json:"comm_no"`
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
	CommNo int64 `json:"comm_no"`
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
	CommNo   int64  `json:"comm_no"`
	NickName string `json:"nick_name"`
	ErrCode  string `json:"err_code"`
	ErrMsg   string `json:"err_msg"`
}

/*
	回复一评论请求
*/
type CommentReplyReq struct {
	CommNo  int64  `json:"comm_no"`
	Title   string `json:"title"`
	Context string `json:"context"`
}

/*
	回复一评论应答
*/
type CommentReplyResp struct {
	NickName string `json:"nick_name"`
	ErrCode  string `json:"err_code"`
	ErrMsg   string `json:"err_msg"`
}

func GetReplyList(userId int64, commNo int64, list *([]comment.Comment)) {
	common.PrintHead("GetReplyList")
	var search comment.Search
	search.ParentCommNo = commNo
	r := comment.New(dbcomm.GetDB(), comment.DEBUG)
	l, _ := r.GetList(search)
	if l == nil {
		return
	}
	for _, v := range l {
		if v.UserId != userId {
			v.IsEableReply = true
		}
		*list = append(*list, v)
		fmt.Println("=======>", v)
		GetReplyList(userId, v.CommNo, list)
	}
	common.PrintTail("GetReplyList")
}

/*
	获取用户发布的评论
*/
func CommentList(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("CommentList")
	userId, _, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
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
	search.ParentCommNo = 0
	r := comment.New(dbcomm.GetDB(), comment.DEBUG)
	l, err := r.GetList(search)
	total, err := r.GetTotal(search)
	for k, v := range l {
		if v.UserId != uId {
			l[k].IsEableReply = true
		}
		var repay_list []comment.Comment
		GetReplyList(uId, v.CommNo, &repay_list)
		l[k].ReplyList = repay_list
	}
	listResp.ErrCode = common.ERR_CODE_SUCCESS
	listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	listResp.CommList = l
	listResp.Total = total
	common.Write_Response(listResp, w, req)
	common.PrintTail("CommentList")
}

/*
	获取用户发布的评论
*/
func MyCommentList(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("CommentList")
	userId, _, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
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
	search.UserId = uId
	r := comment.New(dbcomm.GetDB(), comment.DEBUG)
	l, err := r.GetList(search)
	total, err := r.GetTotal(search)
	listResp.ErrCode = common.ERR_CODE_SUCCESS
	listResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	listResp.CommList = l
	listResp.Total = total
	common.Write_Response(listResp, w, req)
	common.PrintTail("CommentList")
}

/*
	发表评论
*/

func CommentPost(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("CommentPost")
	userId, _, nickName, tokenErr := common.CheckToken(w, req)
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
	commNo := time.Now().UnixNano()
	e.CommNo = commNo
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
	postResp.CommNo = commNo
	postResp.NickName = nickName
	postResp.ErrCode = common.ERR_CODE_SUCCESS
	postResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(postResp, w, req)
	common.PrintTail("CommentPost")
}

/*
	LIKE 某个评论
*/
func CommentLike(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("CommentLike")
	userId, _, _, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)
	fmt.Println("UserId=====>", uId)
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

	//查询用户是否已经点赞
	r := comment_user.New(dbcomm.GetDB(), comment_user.DEBUG)
	var search comment_user.Search
	search.CommNo = likeReq.CommNo
	search.UserId = uId
	search.ActionType = common.COMMENT_LIKE
	_, err = r.Get(search)
	if err == nil {
		likeResp.ErrCode = common.ERR_CODE_EXISTED
		likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED]
		common.Write_Response(likeResp, w, req)
		return
	}

	//记录用户的点赞情况
	tr, err := r.DB.Begin()
	if err != nil {
		likeResp.ErrCode = common.ERR_CODE_DBERROR
		likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(likeResp, w, req)
		return
	}
	rr := comment.New(dbcomm.GetDB(), comment.DEBUG)
	err = rr.UpdateLikes(fmt.Sprintf("%d", likeReq.CommNo), tr)
	if err != nil {
		tr.Rollback()
		likeResp.ErrCode = common.ERR_CODE_DBERROR
		likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(likeResp, w, req)
		return
	}
	rrr := comment_user.New(rr.DB, comment_user.DEBUG)
	var ne comment_user.CommentUser
	ne.UserId = uId
	ne.CommNo = likeReq.CommNo
	ne.ActionType = common.COMMENT_LIKE
	ne.InsertTime = time.Now().Format("2006-01-02 15:04:05")
	err = rrr.InsertEntity(ne, tr)
	if err != nil {
		tr.Rollback()
		likeResp.ErrCode = common.ERR_CODE_DBERROR
		likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(likeResp, w, req)
		return
	}
	tr.Commit()

	//得到这条评论的点赞数量
	var noSearch comment.Search
	noSearch.CommNo = likeReq.CommNo
	re, err := rr.Get(noSearch)
	likeResp.ErrCode = common.ERR_CODE_SUCCESS
	likeResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	likeResp.Likes = re.Likes
	common.Write_Response(likeResp, w, req)
	common.PrintTail("CommentPost")

}

/*
	kILL某个评论
*/
func CommentKill(w http.ResponseWriter, req *http.Request) {
	common.PrintHead("CommentKill")

	userId, _, _, tokenErr := common.CheckToken(w, req)
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

	//查看用户是否投诉
	r := comment_user.New(dbcomm.GetDB(), comment_user.DEBUG)
	var search comment_user.Search
	search.ActionType = common.COMMENT_KILL
	search.CommNo = killReq.CommNo
	search.UserId = uId
	_, err = r.Get(search)
	if err == nil {
		killResp.ErrCode = common.ERR_CODE_EXISTED
		killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_EXISTED]
		common.Write_Response(killResp, w, req)
		return
	}
	//记录用户的投诉
	tr, err := r.DB.Begin()
	if err != nil {
		killResp.ErrCode = common.ERR_CODE_DBERROR
		killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(killResp, w, req)
		return
	}
	rr := comment.New(dbcomm.GetDB(), comment.DEBUG)
	err = rr.UpdateKills(fmt.Sprintf("%d", killReq.CommNo), tr)
	if err != nil {
		tr.Rollback()
		killResp.ErrCode = common.ERR_CODE_DBERROR
		killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(killResp, w, req)
		return
	}
	rrr := comment_user.New(rr.DB, comment_user.DEBUG)
	var ne comment_user.CommentUser
	ne.UserId = uId
	ne.CommNo = killReq.CommNo
	ne.ActionType = common.COMMENT_KILL
	ne.InsertTime = time.Now().Format("2006-01-02 15:04:05")
	err = rrr.InsertEntity(ne, tr)
	if err != nil {
		tr.Rollback()
		killResp.ErrCode = common.ERR_CODE_DBERROR
		killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_DBERROR]
		common.Write_Response(killResp, w, req)
		return
	}
	tr.Commit()

	//得到这条评论的投诉数量
	var noSearch comment.Search
	noSearch.CommNo = killReq.CommNo
	re, err := rr.Get(noSearch)
	killResp.ErrCode = common.ERR_CODE_SUCCESS
	killResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	killResp.Kills = re.Kills

	common.Write_Response(killResp, w, req)
	common.PrintTail("CommentKill")

}

/*
	回复某个评论
*/
func CommentReply(w http.ResponseWriter, req *http.Request) {
	common.PrintTail("CommentReply")
	userId, _, nickName, tokenErr := common.CheckToken(w, req)
	if tokenErr != nil {
		return
	}
	uId, _ := strconv.ParseInt(userId, 10, 64)

	var replyReq CommentReplyReq
	var replyResp CommentReplyResp
	err := json.NewDecoder(req.Body).Decode(&replyReq)
	if err != nil {
		replyResp.ErrCode = common.ERR_CODE_JSONERR
		replyResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_JSONERR] + "请求报文格式有误！" + err.Error()
		common.Write_Response(replyResp, w, req)
		return
	}
	defer req.Body.Close()

	r := comment.New(dbcomm.GetDB(), comment.DEBUG)
	var search comment.Search
	search.CommNo = replyReq.CommNo
	e, err := r.Get(search)
	if err != nil {
		replyResp.ErrCode = common.ERR_CODE_NOTFIND
		replyResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_NOTFIND]
		common.Write_Response(replyResp, w, req)
		return
	}
	var ee comment.Comment
	ee.CommNo = time.Now().UnixNano()
	ee.Title = replyReq.Title
	ee.Context = replyReq.Context
	ee.ParentCommNo = e.CommNo
	ee.UserId = uId
	ee.InsertTime = time.Now().Format("2006-01-02 15:04:05")
	r.InsertEntity(ee, nil)
	replyResp.NickName = nickName
	replyResp.ErrCode = common.ERR_CODE_SUCCESS
	replyResp.ErrMsg = common.ERROR_MAP[common.ERR_CODE_SUCCESS]
	common.Write_Response(replyResp, w, req)
	common.PrintTail("CommentReply")

}
