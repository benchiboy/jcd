package disputes

import (
	"database/sql"
	"fmt"
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
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	MctNo       string `json:"mct_no"`
	MctLoanNo   string `json:"mct_loan_no"`
	Phone       string `json:"phone"`
	Mail        string `json:"mail"`
	DisputeNo   int64  `json:"dispute_no"`
	DisputeMemo string `json:"dispute_memo"`
	Status      string `json:"status"`
	InsertTime  string `json:"insert_time"`
	UpdateTime  string `json:"update_time"`
	UpdateUser  string `json:"update_user"`
	Version     int64  `json:"version"`
	PageNo      int    `json:"page_no"`
	PageSize    int    `json:"page_size"`
	ExtraWhere  string `json:"extra_where"`
	SortFld     string `json:"sort_fld"`
}

type DisputesList struct {
	DB        *sql.DB
	Level     int
	Total     int        `json:"total"`
	Disputess []Disputes `json:"Disputes"`
}

type Disputes struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	MctNo        string `json:"mct_no"`
	MctLoanNo    string `json:"mct_loan_no"`
	Phone        string `json:"phone"`
	Mail         string `json:"mail"`
	DisputeNo    int64  `json:"dispute_no"`
	DisputeMemo  string `json:"dispute_memo"`
	FeedbackMemo string `json:"feedback_memo"`
	Status       string `json:"status"`
	InsertTime   string `json:"insert_time"`
	UpdateTime   string `json:"update_time"`
	UpdateUser   string `json:"update_user"`
	Version      int64  `json:"version"`
}

type Form struct {
	Form Disputes `json:"Disputes"`
}

/*
	说明：创建实例对象
	入参：db:数据库sql.DB, 数据库已经连接, level:日志级别
	出参：实例对象
*/

func New(db *sql.DB, level int) *DisputesList {
	if db == nil {
		log.Println(SQL_SELECT, "Database is nil")
		return nil
	}
	return &DisputesList{DB: db, Total: 0, Disputess: make([]Disputes, 0), Level: level}
}

/*
	说明：创建实例对象
	入参：url:连接数据的url, 数据库还没有CONNECTED, level:日志级别
	出参：实例对象
*/

func NewUrl(url string, level int) *DisputesList {
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
	return &DisputesList{DB: db, Total: 0, Disputess: make([]Disputes, 0), Level: level}
}

/*
	说明：得到符合条件的总条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的总条件, 参数2：如果错误返回错误对象
*/

func (r *DisputesList) GetTotal(s Search) (int, error) {
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

	if s.MctLoanNo != "" {
		where += " and mct_loan_no='" + s.MctLoanNo + "'"
	}

	if s.Phone != "" {
		where += " and phone='" + s.Phone + "'"
	}

	if s.Mail != "" {
		where += " and mail='" + s.Mail + "'"
	}

	if s.DisputeNo != 0 {
		where += " and dispute_no=" + fmt.Sprintf("%d", s.DisputeNo)
	}

	if s.DisputeMemo != "" {
		where += " and dispute_memo='" + s.DisputeMemo + "'"
	}

	if s.Status != "" {
		where += " and status='" + s.Status + "'"
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

	qrySql := fmt.Sprintf("Select count(1) as total from b_disputes   where 1=1 %s", where)
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

func (r DisputesList) Get(s Search) (*Disputes, error) {
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

	if s.MctLoanNo != "" {
		where += " and mct_loan_no='" + s.MctLoanNo + "'"
	}

	if s.Phone != "" {
		where += " and phone='" + s.Phone + "'"
	}

	if s.Mail != "" {
		where += " and mail='" + s.Mail + "'"
	}

	if s.DisputeNo != 0 {
		where += " and dispute_no=" + fmt.Sprintf("%d", s.DisputeNo)
	}

	if s.DisputeMemo != "" {
		where += " and dispute_memo='" + s.DisputeMemo + "'"
	}

	if s.Status != "" {
		where += " and status='" + s.Status + "'"
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

	qrySql := fmt.Sprintf("Select id,user_id,mct_no,mct_loan_no,phone,mail,dispute_no,dispute_memo,status,insert_time,update_time,update_user,version from b_disputes where 1=1 %s ", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p Disputes
	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err := rows.Scan(&p.Id, &p.UserId, &p.MctNo, &p.MctLoanNo, &p.Phone, &p.Mail, &p.DisputeNo, &p.DisputeMemo, &p.Status, &p.InsertTime, &p.UpdateTime, &p.UpdateUser, &p.Version)
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

func (r *DisputesList) GetList(s Search) ([]Disputes, error) {
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

	if s.MctLoanNo != "" {
		where += " and mct_loan_no='" + s.MctLoanNo + "'"
	}

	if s.Phone != "" {
		where += " and phone='" + s.Phone + "'"
	}

	if s.Mail != "" {
		where += " and mail='" + s.Mail + "'"
	}

	if s.DisputeNo != 0 {
		where += " and dispute_no=" + fmt.Sprintf("%d", s.DisputeNo)
	}

	if s.DisputeMemo != "" {
		where += " and dispute_memo='" + s.DisputeMemo + "'"
	}

	if s.Status != "" {
		where += " and status='" + s.Status + "'"
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
		qrySql = fmt.Sprintf("Select id,user_id,mct_no,mct_loan_no,phone,mail,dispute_no,dispute_memo,feedback_memo,status,insert_time,update_time,update_user,version from b_disputes where 1=1 %s", where)
	} else {
		qrySql = fmt.Sprintf("Select id,user_id,mct_no,mct_loan_no,phone,mail,dispute_no,dispute_memo,feedback_memo,status,insert_time,update_time,update_user,version from b_disputes where 1=1 %s Limit %d offset %d", where, s.PageSize, (s.PageNo-1)*s.PageSize)
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

	var p Disputes
	for rows.Next() {
		rows.Scan(&p.Id, &p.UserId, &p.MctNo, &p.MctLoanNo, &p.Phone, &p.Mail, &p.DisputeNo, &p.DisputeMemo, &p.FeedbackMemo, &p.Status, &p.InsertTime, &p.UpdateTime, &p.UpdateUser, &p.Version)
		r.Disputess = append(r.Disputess, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.Disputess, nil
}

/*
	说明：根据主键查询符合条件的记录，并保持成MAP
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象, 参数2：如果错误返回错误对象
*/

func (r *DisputesList) GetExt(s Search) (map[string]string, error) {
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

	if s.MctLoanNo != "" {
		where += " and mct_loan_no='" + s.MctLoanNo + "'"
	}

	if s.Phone != "" {
		where += " and phone='" + s.Phone + "'"
	}

	if s.Mail != "" {
		where += " and mail='" + s.Mail + "'"
	}

	if s.DisputeNo != 0 {
		where += " and dispute_no=" + fmt.Sprintf("%d", s.DisputeNo)
	}

	if s.DisputeMemo != "" {
		where += " and dispute_memo='" + s.DisputeMemo + "'"
	}

	if s.Status != "" {
		where += " and status='" + s.Status + "'"
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

	qrySql := fmt.Sprintf("Select id,user_id,mct_no,mct_loan_no,phone,mail,dispute_no,dispute_memo,status,insert_time,update_time,update_user,version from b_disputes where 1=1 %s ", where)
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

func (r DisputesList) Insert(p Disputes) error {
	l := time.Now()
	exeSql := fmt.Sprintf("Insert into  b_disputes(user_id,mct_no,mct_loan_no,phone,mail,dispute_no,dispute_memo,status,update_user,version)  values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}
	_, err := r.DB.Exec(exeSql, p.UserId, p.MctNo, p.MctLoanNo, p.Phone, p.Mail, p.DisputeNo, p.DisputeMemo, p.Status, p.UpdateUser, p.Version)
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

func (r DisputesList) InsertEntity(p Disputes, tr *sql.Tx) error {
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

	if p.MctLoanNo != "" {
		colNames += "mct_loan_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.MctLoanNo)
	}

	if p.Phone != "" {
		colNames += "phone,"
		colTags += "?,"
		valSlice = append(valSlice, p.Phone)
	}

	if p.Mail != "" {
		colNames += "mail,"
		colTags += "?,"
		valSlice = append(valSlice, p.Mail)
	}

	if p.DisputeNo != 0 {
		colNames += "dispute_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.DisputeNo)
	}

	if p.DisputeMemo != "" {
		colNames += "dispute_memo,"
		colTags += "?,"
		valSlice = append(valSlice, p.DisputeMemo)
	}

	if p.Status != "" {
		colNames += "status,"
		colTags += "?,"
		valSlice = append(valSlice, p.Status)
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

	colNames = strings.TrimRight(colNames, ",")
	colTags = strings.TrimRight(colTags, ",")
	exeSql := fmt.Sprintf("Insert into  b_disputes(%s)  values(%s)", colNames, colTags)
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

func (r DisputesList) InsertMap(m map[string]interface{}, tr *sql.Tx) error {
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

	exeSql := fmt.Sprintf("Insert into  b_disputes(%s)  values(%s)", colNames, colTags)
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

func (r DisputesList) UpdataEntity(keyNo string, p Disputes, tr *sql.Tx) error {
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

	if p.MctLoanNo != "" {
		colNames += "mct_loan_no=?,"

		valSlice = append(valSlice, p.MctLoanNo)
	}

	if p.Phone != "" {
		colNames += "phone=?,"

		valSlice = append(valSlice, p.Phone)
	}

	if p.Mail != "" {
		colNames += "mail=?,"

		valSlice = append(valSlice, p.Mail)
	}

	if p.DisputeNo != 0 {
		colNames += "dispute_no=?,"
		valSlice = append(valSlice, p.DisputeNo)
	}

	if p.DisputeMemo != "" {
		colNames += "dispute_memo=?,"

		valSlice = append(valSlice, p.DisputeMemo)
	}

	if p.Status != "" {
		colNames += "status=?,"

		valSlice = append(valSlice, p.Status)
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

	exeSql := fmt.Sprintf("update  b_disputes  set %s  where id=? ", colNames)
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

func (r DisputesList) UpdateMap(keyNo string, m map[string]interface{}, tr *sql.Tx) error {
	l := time.Now()

	var colNames string
	valSlice := make([]interface{}, 0)
	for k, v := range m {
		colNames += k + "=?,"
		valSlice = append(valSlice, v)
	}
	valSlice = append(valSlice, keyNo)
	colNames = strings.TrimRight(colNames, ",")
	updateSql := fmt.Sprintf("Update b_disputes set %s where id=?", colNames)
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

func (r DisputesList) Delete(keyNo string, tr *sql.Tx) error {
	l := time.Now()
	delSql := fmt.Sprintf("Delete from  b_disputes  where id=?")
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
