package comment

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
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	ParentCommNo int64  `json:"parent_comm_no"`
	CommNo       int64  `json:"comm_no"`
	Title        string `json:"title"`
	Context      string `json:"context"`
	Kills        int64  `json:"kills"`
	Likes        int64  `json:"likes"`
	InsertTime   string `json:"insert_time"`
	UpdateTime   string `json:"update_time"`
	UpdateUser   string `json:"update_user"`
	Version      int64  `json:"version"`
	PageNo       int    `json:"page_no"`
	PageSize     int    `json:"page_size"`
	ExtraWhere   string `json:"extra_where"`
	SortFld      string `json:"sort_fld"`
}

type CommentList struct {
	DB       *sql.DB
	Level    int
	Total    int       `json:"total"`
	Comments []Comment `json:"Comment"`
}

type Comment struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	ParentCommNo int64  `json:"parent_comm_no"`
	CommNo       int64  `json:"comm_no"`
	Title        string `json:"title"`
	Context      string `json:"context"`
	Kills        int64  `json:"kills"`
	Likes        int64  `json:"likes"`
	InsertTime   string `json:"insert_time"`
	UpdateTime   string `json:"update_time"`
	UpdateUser   string `json:"update_user"`
	Version      int64  `json:"version"`
}

type Form struct {
	Form Comment `json:"Comment"`
}

/*
	说明：创建实例对象
	入参：db:数据库sql.DB, 数据库已经连接, level:日志级别
	出参：实例对象
*/

func New(db *sql.DB, level int) *CommentList {
	if db == nil {
		log.Println(SQL_SELECT, "Database is nil")
		return nil
	}
	return &CommentList{DB: db, Total: 0, Comments: make([]Comment, 0), Level: level}
}

/*
	说明：创建实例对象
	入参：url:连接数据的url, 数据库还没有CONNECTED, level:日志级别
	出参：实例对象
*/

func NewUrl(url string, level int) *CommentList {
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
	return &CommentList{DB: db, Total: 0, Comments: make([]Comment, 0), Level: level}
}

/*
	说明：得到符合条件的总条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的总条件, 参数2：如果错误返回错误对象
*/

func (r *CommentList) GetTotal(s Search) (int, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.ParentCommNo != 0 {
		where += " and parent_comm_no=" + fmt.Sprintf("%d", s.ParentCommNo)
	}

	if s.CommNo != 0 {
		where += " and comm_no=" + fmt.Sprintf("%d", s.CommNo)
	}

	if s.Title != "" {
		where += " and title='" + s.Title + "'"
	}

	if s.Context != "" {
		where += " and context='" + s.Context + "'"
	}

	if s.Kills != 0 {
		where += " and kills=" + fmt.Sprintf("%d", s.Kills)
	}

	if s.Likes != 0 {
		where += " and likes=" + fmt.Sprintf("%d", s.Likes)
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

	qrySql := fmt.Sprintf("Select count(1) as total from b_comment   where 1=1 %s", where)
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

func (r CommentList) Get(s Search) (*Comment, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.ParentCommNo != 0 {
		where += " and parent_comm_no=" + fmt.Sprintf("%d", s.ParentCommNo)
	}

	if s.CommNo != 0 {
		where += " and comm_no=" + fmt.Sprintf("%d", s.CommNo)
	}

	if s.Title != "" {
		where += " and title='" + s.Title + "'"
	}

	if s.Context != "" {
		where += " and context='" + s.Context + "'"
	}

	if s.Kills != 0 {
		where += " and kills=" + fmt.Sprintf("%d", s.Kills)
	}

	if s.Likes != 0 {
		where += " and likes=" + fmt.Sprintf("%d", s.Likes)
	}

	if s.InsertTime != "" {
		where += " and insert_time='" + s.InsertTime + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	qrySql := fmt.Sprintf("Select id,user_id,parent_comm_no,comm_no,title,context,kills,likes,update_user,version from b_comment where 1=1 %s ", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p Comment
	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err := rows.Scan(&p.Id, &p.UserId, &p.ParentCommNo, &p.CommNo, &p.Title, &p.Context, &p.Kills, &p.Likes, &p.UpdateUser, &p.Version)
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

func (r *CommentList) GetList(s Search) ([]Comment, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.ParentCommNo != 0 {
		where += " and parent_comm_no=" + fmt.Sprintf("%d", s.ParentCommNo)
	}

	if s.CommNo != 0 {
		where += " and comm_no=" + fmt.Sprintf("%d", s.CommNo)
	}

	if s.Title != "" {
		where += " and title='" + s.Title + "'"
	}

	if s.Context != "" {
		where += " and context='" + s.Context + "'"
	}

	if s.Kills != 0 {
		where += " and kills=" + fmt.Sprintf("%d", s.Kills)
	}

	if s.Likes != 0 {
		where += " and likes=" + fmt.Sprintf("%d", s.Likes)
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
		qrySql = fmt.Sprintf("Select id,user_id,parent_comm_no,comm_no,title,context,kills,likes,insert_time,update_time,update_user,version from b_comment where 1=1 %s", where)
	} else {
		qrySql = fmt.Sprintf("Select id,user_id,parent_comm_no,comm_no,title,context,kills,likes,insert_time,update_time,update_user,version from b_comment where 1=1 %s Limit %d offset %d", where, s.PageSize, (s.PageNo-1)*s.PageSize)
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

	var p Comment
	for rows.Next() {
		rows.Scan(&p.Id, &p.UserId, &p.ParentCommNo, &p.CommNo, &p.Title, &p.Context, &p.Kills, &p.Likes, &p.InsertTime, &p.UpdateTime, &p.UpdateUser, &p.Version)
		r.Comments = append(r.Comments, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.Comments, nil
}

/*
	说明：根据条件查询复核条件对象列表，支持分页查询
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象列表, 参数2：如果错误返回错误对象
*/

func (r *CommentList) GetListExt(s Search, fList []string) ([][]pubtype.Data, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.ParentCommNo != 0 {
		where += " and parent_comm_no=" + fmt.Sprintf("%d", s.ParentCommNo)
	}

	if s.CommNo != 0 {
		where += " and comm_no=" + fmt.Sprintf("%d", s.CommNo)
	}

	if s.Title != "" {
		where += " and title='" + s.Title + "'"
	}

	if s.Context != "" {
		where += " and context='" + s.Context + "'"
	}

	if s.Kills != 0 {
		where += " and kills=" + fmt.Sprintf("%d", s.Kills)
	}

	if s.Likes != 0 {
		where += " and likes=" + fmt.Sprintf("%d", s.Likes)
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
		qrySql = fmt.Sprintf("Select %s from b_comment where 1=1 %s", colNames, where)
	} else {
		qrySql = fmt.Sprintf("Select %s from b_comment where 1=1 %s Limit %d offset %d", colNames, where, s.PageSize, (s.PageNo-1)*s.PageSize)
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

func (r *CommentList) GetExt(s Search) (map[string]string, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.ParentCommNo != 0 {
		where += " and parent_comm_no=" + fmt.Sprintf("%d", s.ParentCommNo)
	}

	if s.CommNo != 0 {
		where += " and comm_no=" + fmt.Sprintf("%d", s.CommNo)
	}

	if s.Title != "" {
		where += " and title='" + s.Title + "'"
	}

	if s.Context != "" {
		where += " and context='" + s.Context + "'"
	}

	if s.Kills != 0 {
		where += " and kills=" + fmt.Sprintf("%d", s.Kills)
	}

	if s.Likes != 0 {
		where += " and likes=" + fmt.Sprintf("%d", s.Likes)
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

	qrySql := fmt.Sprintf("Select id,user_id,parent_comm_no,comm_no,title,context,kills,likes,insert_time,update_time,update_user,version from b_comment where 1=1 %s ", where)
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

func (r CommentList) Insert(p Comment) error {
	l := time.Now()
	exeSql := fmt.Sprintf("Insert into  b_comment(user_id,parent_comm_no,comm_no,title,context,kills,likes,update_user,version)  values(?,?,?,?,?,?,?,?,?,?,?,?)")
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}
	_, err := r.DB.Exec(exeSql, p.UserId, p.ParentCommNo, p.CommNo, p.Title, p.Context, p.Kills, p.Likes, p.UpdateUser, p.Version)
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

func (r CommentList) InsertEntity(p Comment, tr *sql.Tx) error {
	l := time.Now()
	var colNames, colTags string
	valSlice := make([]interface{}, 0)

	if p.UserId != 0 {
		colNames += "user_id,"
		colTags += "?,"
		valSlice = append(valSlice, p.UserId)
	}

	if p.ParentCommNo != 0 {
		colNames += "parent_comm_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.ParentCommNo)
	}

	if p.CommNo != 0 {
		colNames += "comm_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.CommNo)
	}

	if p.Title != "" {
		colNames += "title,"
		colTags += "?,"
		valSlice = append(valSlice, p.Title)
	}

	if p.Context != "" {
		colNames += "context,"
		colTags += "?,"
		valSlice = append(valSlice, p.Context)
	}

	if p.Kills != 0 {
		colNames += "kills,"
		colTags += "?,"
		valSlice = append(valSlice, p.Kills)
	}

	if p.Likes != 0 {
		colNames += "likes,"
		colTags += "?,"
		valSlice = append(valSlice, p.Likes)
	}

	if p.InsertTime != "" {
		colNames += "insert_time,"
		colTags += "?,"
		valSlice = append(valSlice, p.InsertTime)
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
	exeSql := fmt.Sprintf("Insert into  b_comment(%s)  values(%s)", colNames, colTags)
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

func (r CommentList) InsertMap(m map[string]interface{}, tr *sql.Tx) error {
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

	exeSql := fmt.Sprintf("Insert into  b_comment(%s)  values(%s)", colNames, colTags)
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

func (r CommentList) UpdataEntity(keyNo string, p Comment, tr *sql.Tx) error {
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

	if p.ParentCommNo != 0 {
		colNames += "parent_comm_no=?,"
		valSlice = append(valSlice, p.ParentCommNo)
	}

	if p.CommNo != 0 {
		colNames += "comm_no=?,"
		valSlice = append(valSlice, p.CommNo)
	}

	if p.Title != "" {
		colNames += "title=?,"

		valSlice = append(valSlice, p.Title)
	}

	if p.Context != "" {
		colNames += "context=?,"

		valSlice = append(valSlice, p.Context)
	}

	if p.Kills != 0 {
		colNames += "kills=?,"
		valSlice = append(valSlice, p.Kills)
	}

	if p.Likes != 0 {
		colNames += "likes=?,"
		valSlice = append(valSlice, p.Likes)
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

	exeSql := fmt.Sprintf("update  b_comment  set %s  where id=? ", colNames)
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

func (r CommentList) UpdateMap(keyNo string, m map[string]interface{}, tr *sql.Tx) error {
	l := time.Now()

	var colNames string
	valSlice := make([]interface{}, 0)
	for k, v := range m {
		colNames += k + "=?,"
		valSlice = append(valSlice, v)
	}
	valSlice = append(valSlice, keyNo)
	colNames = strings.TrimRight(colNames, ",")
	updateSql := fmt.Sprintf("Update b_comment set %s ,version=version+1 where id=?", colNames)
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

func (r CommentList) Delete(keyNo string, tr *sql.Tx) error {
	l := time.Now()
	delSql := fmt.Sprintf("Delete from  b_comment  where id=?")
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
