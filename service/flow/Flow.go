package flow

import (
	"database/sql"
	"fmt"
	"hcd-gate/service/pubtype"
	"log"
	"strings"
	"time"
)

const (
	SQL_NEWDB   = "NewDB  ===>"
	SQL_INSERT  = "Insert ===>"
	SQL_UPDATE  = "Update ===>"
	SQL_SELECT  = "Select ===>"
	SQL_DELETE  = "Delete ===>"
	SQL_ELAPSED = "Elapsed===>"
	SQL_ERROR   = "Error  ===>"
	SQL_TITLE   = "===================================="
	DEBUG       = 1
	INFO        = 2
)

type Search struct {
	Id         int64   `json:"id"`
	UserId     int64   `json:"user_id"`
	MctNo      string  `json:"mct_no"`
	MctTrxnNo  string  `json:"mct_trxn_no"`
	TrxnNo     int64   `json:"trxn_no"`
	TrxnAmt    int     `json:"trxn_amt"`
	TrxnType   string  `json:"trxn_type"`
	ProcStatus string  `json:"proc_status"`
	ProcMsg    string  `json:"proc_msg"`
	AccountBal float64 `json:"account_bal"`
	CodeUrl    string  `json:"code_url"`
	PrepayId   string  `json:"prepay_id"`
	TrxnMemo   string  `json:"trxn_memo"`
	TrxnDate   string  `json:"trxn_date"`
	DoneDate   string  `json:"done_date"`
	InsertTime string  `json:"insert_time"`
	UpdateTime string  `json:"update_time"`
	UpdateUser string  `json:"update_user"`
	Version    int64   `json:"version"`
	PageNo     int     `json:"page_no"`
	PageSize   int     `json:"page_size"`
	ExtraWhere string  `json:"extra_where"`
	SortFld    string  `json:"sort_fld"`
	CrfUid     string  `json:"crf_uid"`
	UserName   string  `json:"user_name"`
	IdNo       string  `json:"id_no"`
	PhoneNo    string  `json:"phone_no"`
	CardNo     string  `json:"card_no"`
	BankId     string  `json:"bank_id"`
}

type FlowList struct {
	DB    *sql.DB
	Level int
	Total int    `json:"total"`
	Flows []Flow `json:"Flow"`
}

type Flow struct {
	Id         int64   `json:"id"`
	UserId     int64   `json:"user_id"`
	MctNo      string  `json:"mct_no"`
	MctTrxnNo  string  `json:"mct_trxn_no"`
	TrxnNo     int64   `json:"trxn_no"`
	TrxnAmt    int     `json:"trxn_amt"`
	TrxnType   string  `json:"trxn_type"`
	ProcStatus string  `json:"proc_status"`
	ProcMsg    string  `json:"proc_msg"`
	AccountBal float64 `json:"account_bal"`
	CodeUrl    string  `json:"code_url"`
	PrepayId   string  `json:"prepay_id"`
	TrxnMemo   string  `json:"trxn_memo"`
	TrxnDate   string  `json:"trxn_date"`
	DoneDate   string  `json:"done_date"`
	InsertTime string  `json:"insert_time"`
	UpdateTime string  `json:"update_time"`
	UpdateUser string  `json:"update_user"`
	Version    int64   `json:"version"`
	CrfUid     string  `json:"crf_uid"`
	UserName   string  `json:"user_name"`
	IdNo       string  `json:"id_no"`
	PhoneNo    string  `json:"phone_no"`
	CardNo     string  `json:"card_no"`
	BankId     string  `json:"bank_id"`
}

type Form struct {
	Form Flow `json:"Flow"`
}

/*
	说明：创建实例对象
	入参：db:数据库sql.DB, 数据库已经连接, level:日志级别
	出参：实例对象
*/

func New(db *sql.DB, level int) *FlowList {
	if db == nil {
		log.Println(SQL_SELECT, "Database is nil")
		return nil
	}
	return &FlowList{DB: db, Total: 0, Flows: make([]Flow, 0), Level: level}
}

/*
	说明：创建实例对象
	入参：url:连接数据的url, 数据库还没有CONNECTED, level:日志级别
	出参：实例对象
*/

func NewUrl(url string, level int) *FlowList {
	var err error
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Println(SQL_SELECT, "Open database error:", err)
		return nil
	}
	if err = db.Ping(); err != nil {
		log.Println(SQL_SELECT, "Ping database error:", err)
		return nil
	}
	return &FlowList{DB: db, Total: 0, Flows: make([]Flow, 0), Level: level}
}

/*
	说明：得到符合条件的总条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的总条件, 参数2：如果错误返回错误对象
*/

func (r *FlowList) GetTotal(s Search) (int, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.MctNo != "" {
		where += " and mct_no='" + s.MctNo + "'"
	}

	if s.MctTrxnNo != "" {
		where += " and mct_trxn_no='" + s.MctTrxnNo + "'"
	}

	if s.TrxnNo != 0 {
		where += " and trxn_no=" + fmt.Sprintf("%d", s.TrxnNo)
	}

	if s.TrxnAmt != 0 {
		where += " and trxn_amt=" + fmt.Sprintf("%d", s.TrxnAmt)
	}

	if s.TrxnType != "" {
		where += " and trxn_type='" + s.TrxnType + "'"
	}

	if s.ProcStatus != "" {
		where += " and proc_status='" + s.ProcStatus + "'"
	}

	if s.ProcMsg != "" {
		where += " and proc_msg='" + s.ProcMsg + "'"
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.CodeUrl != "" {
		where += " and code_url='" + s.CodeUrl + "'"
	}

	if s.PrepayId != "" {
		where += " and prepay_id='" + s.PrepayId + "'"
	}

	if s.TrxnMemo != "" {
		where += " and trxn_memo='" + s.TrxnMemo + "'"
	}

	if s.TrxnDate != "" {
		where += " and trxn_date='" + s.TrxnDate + "'"
	}

	if s.DoneDate != "" {
		where += " and done_date='" + s.DoneDate + "'"
	}

	if s.InsertTime != "" {
		where += " and insert_time='" + s.InsertTime + "'"
	}

	if s.UpdateTime != "" {
		where += " and update_time='" + s.UpdateTime + "'"
	}

	if s.UpdateUser != "" {
		where += " and update_user='" + s.UpdateUser + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	qrySql := fmt.Sprintf("Select count(1) as total from b_flow   where 1=1 %s", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return 0, err
	}
	defer rows.Close()
	var total int
	for rows.Next() {
		rows.Scan(&total)
	}
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return total, nil
}

/*
	说明：根据主键查询符合条件的条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象, 参数2：如果错误返回错误对象
*/

func (r FlowList) Get(s Search) (*Flow, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.MctNo != "" {
		where += " and mct_no='" + s.MctNo + "'"
	}

	if s.MctTrxnNo != "" {
		where += " and mct_trxn_no='" + s.MctTrxnNo + "'"
	}

	if s.TrxnNo != 0 {
		where += " and trxn_no=" + fmt.Sprintf("%d", s.TrxnNo)
	}

	if s.TrxnAmt != 0 {
		where += " and trxn_amt=" + fmt.Sprintf("%d", s.TrxnAmt)
	}

	if s.TrxnType != "" {
		where += " and trxn_type='" + s.TrxnType + "'"
	}

	if s.ProcStatus != "" {
		where += " and proc_status='" + s.ProcStatus + "'"
	}

	if s.ProcMsg != "" {
		where += " and proc_msg='" + s.ProcMsg + "'"
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.CodeUrl != "" {
		where += " and code_url='" + s.CodeUrl + "'"
	}

	if s.PrepayId != "" {
		where += " and prepay_id='" + s.PrepayId + "'"
	}

	if s.TrxnMemo != "" {
		where += " and trxn_memo='" + s.TrxnMemo + "'"
	}

	if s.TrxnDate != "" {
		where += " and trxn_date='" + s.TrxnDate + "'"
	}

	if s.DoneDate != "" {
		where += " and done_date='" + s.DoneDate + "'"
	}

	if s.InsertTime != "" {
		where += " and insert_time='" + s.InsertTime + "'"
	}

	if s.UpdateTime != "" {
		where += " and update_time='" + s.UpdateTime + "'"
	}

	if s.UpdateUser != "" {
		where += " and update_user='" + s.UpdateUser + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	qrySql := fmt.Sprintf("Select id,user_id,mct_no,mct_trxn_no,trxn_no,trxn_amt,trxn_type,proc_status,proc_msg,account_bal,code_url,prepay_id,trxn_memo,version from b_flow where 1=1 %s order by id desc limit 1", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p Flow
	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err := rows.Scan(&p.Id, &p.UserId, &p.MctNo, &p.MctTrxnNo, &p.TrxnNo, &p.TrxnAmt, &p.TrxnType, &p.ProcStatus, &p.ProcMsg, &p.AccountBal, &p.CodeUrl, &p.PrepayId, &p.TrxnMemo, &p.Version)
		if err != nil {
			log.Println(SQL_ERROR, err.Error())
			return nil, err
		}
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return &p, nil
}

/*
	说明：根据条件查询复核条件对象列表，支持分页查询
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象列表, 参数2：如果错误返回错误对象
*/

func (r *FlowList) GetList(s Search) ([]Flow, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.MctNo != "" {
		where += " and mct_no='" + s.MctNo + "'"
	}

	if s.MctTrxnNo != "" {
		where += " and mct_trxn_no='" + s.MctTrxnNo + "'"
	}

	if s.TrxnNo != 0 {
		where += " and trxn_no=" + fmt.Sprintf("%d", s.TrxnNo)
	}

	if s.TrxnAmt != 0 {
		where += " and trxn_amt=" + fmt.Sprintf("%d", s.TrxnAmt)
	}

	if s.TrxnType != "" {
		where += " and trxn_type='" + s.TrxnType + "'"
	}

	if s.ProcStatus != "" {
		where += " and proc_status='" + s.ProcStatus + "'"
	}

	if s.ProcMsg != "" {
		where += " and proc_msg='" + s.ProcMsg + "'"
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.CodeUrl != "" {
		where += " and code_url='" + s.CodeUrl + "'"
	}

	if s.PrepayId != "" {
		where += " and prepay_id='" + s.PrepayId + "'"
	}

	if s.TrxnMemo != "" {
		where += " and trxn_memo='" + s.TrxnMemo + "'"
	}

	if s.TrxnDate != "" {
		where += " and trxn_date='" + s.TrxnDate + "'"
	}

	if s.DoneDate != "" {
		where += " and done_date='" + s.DoneDate + "'"
	}

	if s.InsertTime != "" {
		where += " and insert_time='" + s.InsertTime + "'"
	}

	if s.UpdateTime != "" {
		where += " and update_time='" + s.UpdateTime + "'"
	}

	if s.UpdateUser != "" {
		where += " and update_user='" + s.UpdateUser + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	var qrySql string
	if s.PageSize == 0 && s.PageNo == 0 {
		qrySql = fmt.Sprintf("Select id,user_id,mct_no,mct_trxn_no,trxn_no,trxn_amt,trxn_type,proc_status,proc_msg,account_bal,code_url,prepay_id,trxn_memo,trxn_date,done_date,insert_time,update_time,update_user,version from b_flow where 1=1 %s", where)
	} else {
		qrySql = fmt.Sprintf("Select id,user_id,mct_no,mct_trxn_no,trxn_no,trxn_amt,trxn_type,proc_status,proc_msg,account_bal,code_url,prepay_id,trxn_memo,trxn_date,done_date,insert_time,update_time,update_user,version from b_flow where 1=1 %s Limit %d offset %d", where, s.PageSize, (s.PageNo-1)*s.PageSize)
	}
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p Flow
	for rows.Next() {
		rows.Scan(&p.Id, &p.UserId, &p.MctNo, &p.MctTrxnNo, &p.TrxnNo, &p.TrxnAmt, &p.TrxnType, &p.ProcStatus, &p.ProcMsg, &p.AccountBal, &p.CodeUrl, &p.PrepayId, &p.TrxnMemo, &p.TrxnDate, &p.DoneDate, &p.InsertTime, &p.UpdateTime, &p.UpdateUser, &p.Version)
		r.Flows = append(r.Flows, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.Flows, nil
}

/*
	说明：根据条件查询复核条件对象列表，支持分页查询
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象列表, 参数2：如果错误返回错误对象
*/

func (r *FlowList) GetListExt(s Search, fList []string) ([][]pubtype.Data, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.MctNo != "" {
		where += " and mct_no='" + s.MctNo + "'"
	}

	if s.MctTrxnNo != "" {
		where += " and mct_trxn_no='" + s.MctTrxnNo + "'"
	}

	if s.TrxnNo != 0 {
		where += " and trxn_no=" + fmt.Sprintf("%d", s.TrxnNo)
	}

	if s.TrxnAmt != 0 {
		where += " and trxn_amt=" + fmt.Sprintf("%d", s.TrxnAmt)
	}

	if s.TrxnType != "" {
		where += " and trxn_type='" + s.TrxnType + "'"
	}

	if s.ProcStatus != "" {
		where += " and proc_status='" + s.ProcStatus + "'"
	}

	if s.ProcMsg != "" {
		where += " and proc_msg='" + s.ProcMsg + "'"
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.CodeUrl != "" {
		where += " and code_url='" + s.CodeUrl + "'"
	}

	if s.PrepayId != "" {
		where += " and prepay_id='" + s.PrepayId + "'"
	}

	if s.TrxnMemo != "" {
		where += " and trxn_memo='" + s.TrxnMemo + "'"
	}

	if s.TrxnDate != "" {
		where += " and trxn_date='" + s.TrxnDate + "'"
	}

	if s.DoneDate != "" {
		where += " and done_date='" + s.DoneDate + "'"
	}

	if s.InsertTime != "" {
		where += " and insert_time='" + s.InsertTime + "'"
	}

	if s.UpdateTime != "" {
		where += " and update_time='" + s.UpdateTime + "'"
	}

	if s.UpdateUser != "" {
		where += " and update_user='" + s.UpdateUser + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	colNames := ""
	for _, v := range fList {
		colNames += v + ","

	}
	colNames = strings.TrimRight(colNames, ",")

	var qrySql string
	if s.PageSize == 0 && s.PageNo == 0 {
		qrySql = fmt.Sprintf("Select %s from b_flow where 1=1 %s", colNames, where)
	} else {
		qrySql = fmt.Sprintf("Select %s from b_flow where 1=1 %s Limit %d offset %d", colNames, where, s.PageSize, (s.PageNo-1)*s.PageSize)
	}
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	Columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(Columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	rowData := make([][]pubtype.Data, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		colData := make([]pubtype.Data, 0)
		for k, _ := range values {
			d := new(pubtype.Data)
			d.FieldName = Columns[k]
			d.FieldValue = string(values[k])
			colData = append(colData, *d)
		}
		//extra flow_batch_id
		d2 := new(pubtype.Data)
		d2.FieldName = "flow_batch_id"
		d2.FieldValue = string(values[0])
		colData = append(colData, *d2)

		rowData = append(rowData, colData)
	}

	log.Println(SQL_ELAPSED, "==========>>>>>>>>>>>", rowData)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return rowData, nil
}

/*
	说明：根据主键查询符合条件的记录，并保持成MAP
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象, 参数2：如果错误返回错误对象
*/

func (r *FlowList) GetExt(s Search) (map[string]string, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.MctNo != "" {
		where += " and mct_no='" + s.MctNo + "'"
	}

	if s.MctTrxnNo != "" {
		where += " and mct_trxn_no='" + s.MctTrxnNo + "'"
	}

	if s.TrxnNo != 0 {
		where += " and trxn_no=" + fmt.Sprintf("%d", s.TrxnNo)
	}

	if s.TrxnAmt != 0 {
		where += " and trxn_amt=" + fmt.Sprintf("%d", s.TrxnAmt)
	}

	if s.TrxnType != "" {
		where += " and trxn_type='" + s.TrxnType + "'"
	}

	if s.ProcStatus != "" {
		where += " and proc_status='" + s.ProcStatus + "'"
	}

	if s.ProcMsg != "" {
		where += " and proc_msg='" + s.ProcMsg + "'"
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.CodeUrl != "" {
		where += " and code_url='" + s.CodeUrl + "'"
	}

	if s.PrepayId != "" {
		where += " and prepay_id='" + s.PrepayId + "'"
	}

	if s.TrxnMemo != "" {
		where += " and trxn_memo='" + s.TrxnMemo + "'"
	}

	if s.TrxnDate != "" {
		where += " and trxn_date='" + s.TrxnDate + "'"
	}

	if s.DoneDate != "" {
		where += " and done_date='" + s.DoneDate + "'"
	}

	if s.InsertTime != "" {
		where += " and insert_time='" + s.InsertTime + "'"
	}

	if s.UpdateTime != "" {
		where += " and update_time='" + s.UpdateTime + "'"
	}

	if s.UpdateUser != "" {
		where += " and update_user='" + s.UpdateUser + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	qrySql := fmt.Sprintf("Select id,user_id,mct_no,mct_trxn_no,trxn_no,trxn_amt,trxn_type,proc_status,proc_msg,account_bal,code_url,prepay_id,trxn_memo,trxn_date,done_date,insert_time,update_time,update_user,version from b_flow where 1=1 %s ", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	Columns, _ := rows.Columns()

	values := make([]sql.RawBytes, len(Columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err = rows.Scan(scanArgs...)
	}

	fldValMap := make(map[string]string)
	for k, v := range Columns {
		fldValMap[v] = string(values[k])
	}

	log.Println(SQL_ELAPSED, "==========>>>>>>>>>>>", fldValMap)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return fldValMap, nil

}

/*
	说明：插入对象到数据表中，这个方法要求对象的各个属性必须赋值
	入参：p:插入的对象
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r FlowList) Insert(p Flow) error {
	l := time.Now()
	exeSql := fmt.Sprintf("Insert into  b_flow(user_id,mct_no,mct_trxn_no,trxn_no,trxn_amt,trxn_type,proc_status,proc_msg,account_bal,code_url,prepay_id,trxn_memo,trxn_date,done_date,update_user,version)  values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}
	_, err := r.DB.Exec(exeSql, p.UserId, p.MctNo, p.MctTrxnNo, p.TrxnNo, p.TrxnAmt, p.TrxnType, p.ProcStatus, p.ProcMsg, p.AccountBal, p.CodeUrl, p.PrepayId, p.TrxnMemo, p.TrxnDate, p.DoneDate, p.UpdateUser, p.Version)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}

/*
	说明：插入对象到数据表中，这个方法会判读对象的各个属性，如果属性不为空，才加入插入列中；
	入参：p:插入的对象
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r FlowList) InsertEntity(p Flow, tr *sql.Tx) error {
	l := time.Now()
	var colNames, colTags string
	valSlice := make([]interface{}, 0)

	if p.UserId != 0 {
		colNames += "user_id,"
		colTags += "?,"
		valSlice = append(valSlice, p.UserId)
	}

	if p.MctNo != "" {
		colNames += "mct_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.MctNo)
	}

	if p.MctTrxnNo != "" {
		colNames += "mct_trxn_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.MctTrxnNo)
	}

	if p.TrxnNo != 0 {
		colNames += "trxn_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.TrxnNo)
	}

	if p.TrxnAmt != 0 {
		colNames += "trxn_amt,"
		colTags += "?,"
		valSlice = append(valSlice, p.TrxnAmt)
	}

	if p.TrxnType != "" {
		colNames += "trxn_type,"
		colTags += "?,"
		valSlice = append(valSlice, p.TrxnType)
	}

	if p.ProcStatus != "" {
		colNames += "proc_status,"
		colTags += "?,"
		valSlice = append(valSlice, p.ProcStatus)
	}

	if p.ProcMsg != "" {
		colNames += "proc_msg,"
		colTags += "?,"
		valSlice = append(valSlice, p.ProcMsg)
	}

	if p.AccountBal != 0.00 {
		colNames += "account_bal,"
		colTags += "?,"
		valSlice = append(valSlice, p.AccountBal)
	}

	if p.CodeUrl != "" {
		colNames += "code_url,"
		colTags += "?,"
		valSlice = append(valSlice, p.CodeUrl)
	}

	if p.PrepayId != "" {
		colNames += "prepay_id,"
		colTags += "?,"
		valSlice = append(valSlice, p.PrepayId)
	}

	if p.TrxnMemo != "" {
		colNames += "trxn_memo,"
		colTags += "?,"
		valSlice = append(valSlice, p.TrxnMemo)
	}

	if p.TrxnDate != "" {
		colNames += "trxn_date,"
		colTags += "?,"
		valSlice = append(valSlice, p.TrxnDate)
	}

	if p.DoneDate != "" {
		colNames += "done_date,"
		colTags += "?,"
		valSlice = append(valSlice, p.DoneDate)
	}

	if p.UpdateUser != "" {
		colNames += "update_user,"
		colTags += "?,"
		valSlice = append(valSlice, p.UpdateUser)
	}

	if p.Version != 0 {
		colNames += "version,"
		colTags += "?,"
		valSlice = append(valSlice, p.Version)
	}
	if p.CrfUid != "" {
		colNames += "crf_uid,"
		colTags += "?,"
		valSlice = append(valSlice, p.CrfUid)
	}

	if p.UserName != "" {
		colNames += "user_name,"
		colTags += "?,"
		valSlice = append(valSlice, p.UserName)
	}
	if p.IdNo != "" {
		colNames += "id_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.IdNo)
	}

	if p.CardNo != "" {
		colNames += "card_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.CardNo)
	}

	if p.PhoneNo != "" {
		colNames += "phone_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.PhoneNo)
	}

	if p.BankId != "" {
		colNames += "bank_id,"
		colTags += "?,"
		valSlice = append(valSlice, p.BankId)
	}

	colNames = strings.TrimRight(colNames, ",")
	colTags = strings.TrimRight(colTags, ",")
	exeSql := fmt.Sprintf("Insert into  b_flow(%s)  values(%s)", colNames, colTags)
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}

	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(exeSql)
	} else {
		stmt, err = tr.Prepare(exeSql)
	}
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(valSlice...)
	if err != nil {
		log.Println(SQL_INSERT, "Insert data error: %v\n", err)
		return err
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_INSERT, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_INSERT, "RowsAffected:", RowsAffected)
	}

	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}

/*
	说明：插入一个MAP到数据表中；
	入参：m:插入的Map
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r FlowList) InsertMap(m map[string]interface{}, tr *sql.Tx) error {
	l := time.Now()
	var colNames, colTags string
	valSlice := make([]interface{}, 0)
	for k, v := range m {
		colNames += k + ","
		colTags += "?,"
		valSlice = append(valSlice, v)
	}
	colNames = strings.TrimRight(colNames, ",")
	colTags = strings.TrimRight(colTags, ",")

	exeSql := fmt.Sprintf("Insert into  b_flow(%s)  values(%s)", colNames, colTags)
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}

	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(exeSql)
	} else {
		stmt, err = tr.Prepare(exeSql)
	}

	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(valSlice...)
	if err != nil {
		log.Println(SQL_INSERT, "insert data error: %v\n", err)
		return err
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_INSERT, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_INSERT, "RowsAffected:", RowsAffected)
	}

	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}

/*
	说明：插入对象到数据表中，这个方法会判读对象的各个属性，如果属性不为空，才加入插入列中；
	入参：p:插入的对象
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r FlowList) UpdataEntity(keyNo string, p Flow, tr *sql.Tx) error {
	l := time.Now()
	var colNames string
	valSlice := make([]interface{}, 0)

	if p.Id != 0 {
		colNames += "id=?,"
		valSlice = append(valSlice, p.Id)
	}

	if p.UserId != 0 {
		colNames += "user_id=?,"
		valSlice = append(valSlice, p.UserId)
	}

	if p.MctNo != "" {
		colNames += "mct_no=?,"

		valSlice = append(valSlice, p.MctNo)
	}

	if p.MctTrxnNo != "" {
		colNames += "mct_trxn_no=?,"

		valSlice = append(valSlice, p.MctTrxnNo)
	}

	if p.TrxnNo != 0 {
		colNames += "trxn_no=?,"
		valSlice = append(valSlice, p.TrxnNo)
	}

	if p.TrxnAmt != 0 {
		colNames += "trxn_amt=?,"
		valSlice = append(valSlice, p.TrxnAmt)
	}

	if p.TrxnType != "" {
		colNames += "trxn_type=?,"

		valSlice = append(valSlice, p.TrxnType)
	}

	if p.ProcStatus != "" {
		colNames += "proc_status=?,"

		valSlice = append(valSlice, p.ProcStatus)
	}

	if p.ProcMsg != "" {
		colNames += "proc_msg=?,"

		valSlice = append(valSlice, p.ProcMsg)
	}

	if p.AccountBal != 0.00 {
		colNames += "account_bal=?,"
		valSlice = append(valSlice, p.AccountBal)
	}

	if p.CodeUrl != "" {
		colNames += "code_url=?,"

		valSlice = append(valSlice, p.CodeUrl)
	}

	if p.PrepayId != "" {
		colNames += "prepay_id=?,"

		valSlice = append(valSlice, p.PrepayId)
	}

	if p.TrxnMemo != "" {
		colNames += "trxn_memo=?,"

		valSlice = append(valSlice, p.TrxnMemo)
	}

	if p.TrxnDate != "" {
		colNames += "trxn_date=?,"

		valSlice = append(valSlice, p.TrxnDate)
	}

	if p.DoneDate != "" {
		colNames += "done_date=?,"

		valSlice = append(valSlice, p.DoneDate)
	}

	if p.InsertTime != "" {
		colNames += "insert_time=?,"

		valSlice = append(valSlice, p.InsertTime)
	}

	if p.UpdateTime != "" {
		colNames += "update_time=?,"

		valSlice = append(valSlice, p.UpdateTime)
	}

	if p.UpdateUser != "" {
		colNames += "update_user=?,"

		valSlice = append(valSlice, p.UpdateUser)
	}

	if p.Version != 0 {
		colNames += "version=?,"
		valSlice = append(valSlice, p.Version)
	}

	colNames = strings.TrimRight(colNames, ",")
	valSlice = append(valSlice, keyNo)

	exeSql := fmt.Sprintf("update  b_flow  set %s  where id=? ", colNames)
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}

	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(exeSql)
	} else {
		stmt, err = tr.Prepare(exeSql)
	}

	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(valSlice...)
	if err != nil {
		log.Println(SQL_INSERT, "Update data error: %v\n", err)
		return err
	}
	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_INSERT, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_INSERT, "RowsAffected:", RowsAffected)
	}

	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}

/*
	说明：根据更新主键及更新Map值更新数据表；
	入参：keyNo:更新数据的关键条件，m:更新数据列的Map
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r FlowList) UpdateMap(keyNo string, m map[string]interface{}, tr *sql.Tx) error {
	l := time.Now()

	var colNames string
	valSlice := make([]interface{}, 0)
	for k, v := range m {
		colNames += k + "=?,"
		valSlice = append(valSlice, v)
	}
	valSlice = append(valSlice, keyNo)
	colNames = strings.TrimRight(colNames, ",")
	updateSql := fmt.Sprintf("Update b_flow set %s where trxn_no=?", colNames)
	if r.Level == DEBUG {
		log.Println(SQL_UPDATE, updateSql)
	}
	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(updateSql)
	} else {
		stmt, err = tr.Prepare(updateSql)
	}

	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	ret, err := stmt.Exec(valSlice...)
	if err != nil {
		log.Println(SQL_UPDATE, "Update data error: %v\n", err)
		return err
	}
	defer stmt.Close()

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_UPDATE, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_UPDATE, "RowsAffected:", RowsAffected)
	}
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}

/*
	说明：根据主键删除一条数据；
	入参：keyNo:要删除的主键值
	出参：参数1：如果出错，返回错误对象；成功返回nil
*/

func (r FlowList) Delete(keyNo string, tr *sql.Tx) error {
	l := time.Now()
	delSql := fmt.Sprintf("Delete from  b_flow  where id=?")
	if r.Level == DEBUG {
		log.Println(SQL_UPDATE, delSql)
	}

	var stmt *sql.Stmt
	var err error
	if tr == nil {
		stmt, err = r.DB.Prepare(delSql)
	} else {
		stmt, err = tr.Prepare(delSql)
	}

	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return err
	}
	ret, err := stmt.Exec(keyNo)
	if err != nil {
		log.Println(SQL_DELETE, "Delete error: %v\n", err)
		return err
	}
	defer stmt.Close()

	if LastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println(SQL_DELETE, "LastInsertId:", LastInsertId)
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println(SQL_DELETE, "RowsAffected:", RowsAffected)
	}
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return nil
}
