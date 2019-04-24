package account

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
	Id            int64  `json:"id"`
	UserId        int64  `json:"user_id"`
	WxOpenId      string `json:"wx_open_id"`
	WxUnionId     string `json:"wx_union_id"`
	WxSessionKey  string `json:"wx_session_key"`
	LoginMode     int64  `json:"login_mode"`
	LoginName     string `json:"login_name"`
	LoginPass     string `json:"login_pass"`
	Status        int64  `json:"status"`
	AvatarUrl     string `json:"avatar_url"`
	NickName      string `json:"nick_name"`
	Gender        int64  `json:"gender"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	Language      string `json:"language"`
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`

	Errors      int64   `json:"errors"`
	AccountBal  float64 `json:"account_bal"`
	Market      string  `json:"market"`
	RandomNo    int64   `json:"random_no"`
	CreatedTime string  `json:"created_time"`
	UpdatedTime string  `json:"updated_time"`
	Memo        string  `json:"memo"`
	Version     int64   `json:"version"`
	PageNo      int     `json:"page_no"`
	PageSize    int     `json:"page_size"`
	ExtraWhere  string  `json:"extra_where"`
	SortFld     string  `json:"sort_fld"`
}

type AccountList struct {
	DB       *sql.DB
	Level    int
	Total    int       `json:"total"`
	Accounts []Account `json:"Account"`
}

type Account struct {
	Id            int64   `json:"id"`
	UserId        int64   `json:"user_id"`
	WxOpenId      string  `json:"wx_open_id"`
	WxUnionId     string  `json:"wx_union_id"`
	WxSessionKey  string  `json:"wx_session_key"`
	LoginMode     int64   `json:"login_mode"`
	LoginName     string  `json:"login_name"`
	LoginPass     string  `json:"login_pass"`
	Status        int64   `json:"status"`
	AvatarUrl     string  `json:"avatar_url"`
	NickName      string  `json:"nick_name"`
	Gender        int     `json:"gender"`
	City          string  `json:"city"`
	Province      string  `json:"province"`
	Country       string  `json:"country"`
	Language      string  `json:"language"`
	EncryptedData string  `json:"encryptedData"`
	Iv            string  `json:"iv"`
	Errors        int64   `json:"errors"`
	AccountBal    float64 `json:"account_bal"`
	Market        string  `json:"market"`
	RandomNo      int64   `json:"random_no"`
	CreatedTime   string  `json:"created_time"`
	UpdatedTime   string  `json:"updated_time"`
	Memo          string  `json:"memo"`
	Version       int64   `json:"version"`
}

type Form struct {
	Form Account `json:"Account"`
}

/*
	说明：创建实例对象
	入参：db:数据库sql.DB, 数据库已经连接, level:日志级别
	出参：实例对象
*/

func New(db *sql.DB, level int) *AccountList {
	if db == nil {
		log.Println(SQL_SELECT, "Database is nil")
		return nil
	}
	return &AccountList{DB: db, Total: 0, Accounts: make([]Account, 0), Level: level}
}

/*
	说明：创建实例对象
	入参：url:连接数据的url, 数据库还没有CONNECTED, level:日志级别
	出参：实例对象
*/

func NewUrl(url string, level int) *AccountList {
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
	return &AccountList{DB: db, Total: 0, Accounts: make([]Account, 0), Level: level}
}

/*
	说明：得到符合条件的总条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的总条件, 参数2：如果错误返回错误对象
*/

func (r *AccountList) GetTotal(s Search) (int, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.WxOpenId != "" {
		where += " and wx_open_id='" + s.WxOpenId + "'"
	}

	if s.WxUnionId != "" {
		where += " and wx_union_id='" + s.WxUnionId + "'"
	}

	if s.WxSessionKey != "" {
		where += " and wx_session_key='" + s.WxSessionKey + "'"
	}

	if s.LoginMode != 0 {
		where += " and login_mode=" + fmt.Sprintf("%d", s.LoginMode)
	}

	if s.LoginName != "" {
		where += " and login_name='" + s.LoginName + "'"
	}

	if s.LoginPass != "" {
		where += " and login_pass='" + s.LoginPass + "'"
	}

	if s.Status != 0 {
		where += " and status=" + fmt.Sprintf("%d", s.Status)
	}

	if s.AvatarUrl != "" {
		where += " and avatar_url='" + s.AvatarUrl + "'"
	}

	if s.NickName != "" {
		where += " and nick_name='" + s.NickName + "'"
	}

	if s.Gender != 0 {
		where += " and gender=" + fmt.Sprintf("%d", s.Gender)
	}

	if s.City != "" {
		where += " and city='" + s.City + "'"
	}

	if s.Province != "" {
		where += " and province='" + s.Province + "'"
	}

	if s.Country != "" {
		where += " and country='" + s.Country + "'"
	}

	if s.Language != "" {
		where += " and language='" + s.Language + "'"
	}

	if s.Errors != 0 {
		where += " and errors=" + fmt.Sprintf("%d", s.Errors)
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.Market != "" {
		where += " and market='" + s.Market + "'"
	}

	if s.RandomNo != 0 {
		where += " and random_no=" + fmt.Sprintf("%d", s.RandomNo)
	}

	if s.CreatedTime != "" {
		where += " and created_time='" + s.CreatedTime + "'"
	}

	if s.UpdatedTime != "" {
		where += " and updated_time='" + s.UpdatedTime + "'"
	}

	if s.Memo != "" {
		where += " and memo='" + s.Memo + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	qrySql := fmt.Sprintf("Select count(1) as total from b_account   where 1=1 %s", where)
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

func (r AccountList) Get(s Search) (*Account, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.WxOpenId != "" {
		where += " and wx_open_id='" + s.WxOpenId + "'"
	}

	if s.WxUnionId != "" {
		where += " and wx_union_id='" + s.WxUnionId + "'"
	}

	if s.WxSessionKey != "" {
		where += " and wx_session_key='" + s.WxSessionKey + "'"
	}

	if s.LoginMode != 0 {
		where += " and login_mode=" + fmt.Sprintf("%d", s.LoginMode)
	}

	if s.LoginName != "" {
		where += " and login_name='" + s.LoginName + "'"
	}

	if s.LoginPass != "" {
		where += " and login_pass='" + s.LoginPass + "'"
	}

	if s.Status != 0 {
		where += " and status=" + fmt.Sprintf("%d", s.Status)
	}

	if s.AvatarUrl != "" {
		where += " and avatar_url='" + s.AvatarUrl + "'"
	}

	if s.NickName != "" {
		where += " and nick_name='" + s.NickName + "'"
	}

	if s.Gender != 0 {
		where += " and gender=" + fmt.Sprintf("%d", s.Gender)
	}

	if s.City != "" {
		where += " and city='" + s.City + "'"
	}

	if s.Province != "" {
		where += " and province='" + s.Province + "'"
	}

	if s.Country != "" {
		where += " and country='" + s.Country + "'"
	}

	if s.Language != "" {
		where += " and language='" + s.Language + "'"
	}

	if s.Errors != 0 {
		where += " and errors=" + fmt.Sprintf("%d", s.Errors)
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.Market != "" {
		where += " and market='" + s.Market + "'"
	}

	if s.RandomNo != 0 {
		where += " and random_no=" + fmt.Sprintf("%d", s.RandomNo)
	}

	if s.CreatedTime != "" {
		where += " and created_time='" + s.CreatedTime + "'"
	}

	if s.UpdatedTime != "" {
		where += " and updated_time='" + s.UpdatedTime + "'"
	}

	if s.Memo != "" {
		where += " and memo='" + s.Memo + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	qrySql := fmt.Sprintf("Select id,user_id,wx_open_id,wx_union_id,wx_session_key,login_mode,login_name,login_pass,status,avatar_url,nick_name,gender,city,province,country,language,errors,account_bal,market,random_no,memo,version from b_account where 1=1 %s ", where)
	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p Account
	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err := rows.Scan(&p.Id, &p.UserId, &p.WxOpenId, &p.WxUnionId, &p.WxSessionKey, &p.LoginMode, &p.LoginName, &p.LoginPass, &p.Status, &p.AvatarUrl, &p.NickName, &p.Gender, &p.City, &p.Province, &p.Country, &p.Language, &p.Errors, &p.AccountBal, &p.Market, &p.RandomNo, &p.Memo, &p.Version)
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

func (r *AccountList) GetList(s Search) ([]Account, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.WxOpenId != "" {
		where += " and wx_open_id='" + s.WxOpenId + "'"
	}

	if s.WxUnionId != "" {
		where += " and wx_union_id='" + s.WxUnionId + "'"
	}

	if s.WxSessionKey != "" {
		where += " and wx_session_key='" + s.WxSessionKey + "'"
	}

	if s.LoginMode != 0 {
		where += " and login_mode=" + fmt.Sprintf("%d", s.LoginMode)
	}

	if s.LoginName != "" {
		where += " and login_name='" + s.LoginName + "'"
	}

	if s.LoginPass != "" {
		where += " and login_pass='" + s.LoginPass + "'"
	}

	if s.Status != 0 {
		where += " and status=" + fmt.Sprintf("%d", s.Status)
	}

	if s.AvatarUrl != "" {
		where += " and avatar_url='" + s.AvatarUrl + "'"
	}

	if s.NickName != "" {
		where += " and nick_name='" + s.NickName + "'"
	}

	if s.Gender != 0 {
		where += " and gender=" + fmt.Sprintf("%d", s.Gender)
	}

	if s.City != "" {
		where += " and city='" + s.City + "'"
	}

	if s.Province != "" {
		where += " and province='" + s.Province + "'"
	}

	if s.Country != "" {
		where += " and country='" + s.Country + "'"
	}

	if s.Language != "" {
		where += " and language='" + s.Language + "'"
	}

	if s.Errors != 0 {
		where += " and errors=" + fmt.Sprintf("%d", s.Errors)
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.Market != "" {
		where += " and market='" + s.Market + "'"
	}

	if s.RandomNo != 0 {
		where += " and random_no=" + fmt.Sprintf("%d", s.RandomNo)
	}

	if s.CreatedTime != "" {
		where += " and created_time='" + s.CreatedTime + "'"
	}

	if s.UpdatedTime != "" {
		where += " and updated_time='" + s.UpdatedTime + "'"
	}

	if s.Memo != "" {
		where += " and memo='" + s.Memo + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	var qrySql string
	if s.PageSize == 0 && s.PageNo == 0 {
		qrySql = fmt.Sprintf("Select id,user_id,wx_open_id,wx_union_id,wx_session_key,login_mode,login_name,login_pass,status,avatar_url,nick_name,gender,city,province,country,language,errors,account_bal,market,random_no,created_time,updated_time,memo,version from b_account where 1=1 %s", where)
	} else {
		qrySql = fmt.Sprintf("Select id,user_id,wx_open_id,wx_union_id,wx_session_key,login_mode,login_name,login_pass,status,avatar_url,nick_name,gender,city,province,country,language,errors,account_bal,market,random_no,created_time,updated_time,memo,version from b_account where 1=1 %s Limit %d offset %d", where, s.PageSize, (s.PageNo-1)*s.PageSize)
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

	var p Account
	for rows.Next() {
		rows.Scan(&p.Id, &p.UserId, &p.WxOpenId, &p.WxUnionId, &p.WxSessionKey, &p.LoginMode, &p.LoginName, &p.LoginPass, &p.Status, &p.AvatarUrl, &p.NickName, &p.Gender, &p.City, &p.Province, &p.Country, &p.Language, &p.Errors, &p.AccountBal, &p.Market, &p.RandomNo, &p.CreatedTime, &p.UpdatedTime, &p.Memo, &p.Version)
		r.Accounts = append(r.Accounts, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.Accounts, nil
}

/*
	说明：根据主键查询符合条件的记录，并保持成MAP
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象, 参数2：如果错误返回错误对象
*/

func (r *AccountList) GetExt(s Search) (map[string]string, error) {
	var where string
	l := time.Now()

	if s.Id != 0 {
		where += " and id=" + fmt.Sprintf("%d", s.Id)
	}

	if s.UserId != 0 {
		where += " and user_id=" + fmt.Sprintf("%d", s.UserId)
	}

	if s.WxOpenId != "" {
		where += " and wx_open_id='" + s.WxOpenId + "'"
	}

	if s.WxUnionId != "" {
		where += " and wx_union_id='" + s.WxUnionId + "'"
	}

	if s.WxSessionKey != "" {
		where += " and wx_session_key='" + s.WxSessionKey + "'"
	}

	if s.LoginMode != 0 {
		where += " and login_mode=" + fmt.Sprintf("%d", s.LoginMode)
	}

	if s.LoginName != "" {
		where += " and login_name='" + s.LoginName + "'"
	}

	if s.LoginPass != "" {
		where += " and login_pass='" + s.LoginPass + "'"
	}

	if s.Status != 0 {
		where += " and status=" + fmt.Sprintf("%d", s.Status)
	}

	if s.AvatarUrl != "" {
		where += " and avatar_url='" + s.AvatarUrl + "'"
	}

	if s.NickName != "" {
		where += " and nick_name='" + s.NickName + "'"
	}

	if s.Gender != 0 {
		where += " and gender=" + fmt.Sprintf("%d", s.Gender)
	}

	if s.City != "" {
		where += " and city='" + s.City + "'"
	}

	if s.Province != "" {
		where += " and province='" + s.Province + "'"
	}

	if s.Country != "" {
		where += " and country='" + s.Country + "'"
	}

	if s.Language != "" {
		where += " and language='" + s.Language + "'"
	}

	if s.Errors != 0 {
		where += " and errors=" + fmt.Sprintf("%d", s.Errors)
	}

	if s.AccountBal != 0 {
		where += " and account_bal=" + fmt.Sprintf("%f", s.AccountBal)
	}

	if s.Market != "" {
		where += " and market='" + s.Market + "'"
	}

	if s.RandomNo != 0 {
		where += " and random_no=" + fmt.Sprintf("%d", s.RandomNo)
	}

	if s.CreatedTime != "" {
		where += " and created_time='" + s.CreatedTime + "'"
	}

	if s.UpdatedTime != "" {
		where += " and updated_time='" + s.UpdatedTime + "'"
	}

	if s.Memo != "" {
		where += " and memo='" + s.Memo + "'"
	}

	if s.Version != 0 {
		where += " and version=" + fmt.Sprintf("%d", s.Version)
	}

	qrySql := fmt.Sprintf("Select id,user_id,wx_open_id,wx_union_id,wx_session_key,login_mode,login_name,login_pass,status,avatar_url,nick_name,gender,city,province,country,language,errors,account_bal,market,random_no,created_time,updated_time,memo,version from b_account where 1=1 %s ", where)
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

func (r AccountList) Insert(p Account) error {
	l := time.Now()
	exeSql := fmt.Sprintf("Insert into  b_account(user_id,wx_open_id,wx_union_id,wx_session_key,login_mode,login_name,login_pass,status,avatar_url,nick_name,gender,city,province,country,language,errors,account_bal,market,random_no,created_time,updated_time,memo,version)  values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if r.Level == DEBUG {
		log.Println(SQL_INSERT, exeSql)
	}
	_, err := r.DB.Exec(exeSql, p.UserId, p.WxOpenId, p.WxUnionId, p.WxSessionKey, p.LoginMode, p.LoginName, p.LoginPass, p.Status, p.AvatarUrl, p.NickName, p.Gender, p.City, p.Province, p.Country, p.Language, p.Errors, p.AccountBal, p.Market, p.RandomNo, p.CreatedTime, p.UpdatedTime, p.Memo, p.Version)
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

func (r AccountList) InsertEntity(p Account, tr *sql.Tx) error {
	l := time.Now()
	var colNames, colTags string
	valSlice := make([]interface{}, 0)

	if p.UserId != 0 {
		colNames += "user_id,"
		colTags += "?,"
		valSlice = append(valSlice, p.UserId)
	}

	if p.WxOpenId != "" {
		colNames += "wx_open_id,"
		colTags += "?,"
		valSlice = append(valSlice, p.WxOpenId)
	}

	if p.WxUnionId != "" {
		colNames += "wx_union_id,"
		colTags += "?,"
		valSlice = append(valSlice, p.WxUnionId)
	}

	if p.WxSessionKey != "" {
		colNames += "wx_session_key,"
		colTags += "?,"
		valSlice = append(valSlice, p.WxSessionKey)
	}

	if p.LoginMode != 0 {
		colNames += "login_mode,"
		colTags += "?,"
		valSlice = append(valSlice, p.LoginMode)
	}

	if p.LoginName != "" {
		colNames += "login_name,"
		colTags += "?,"
		valSlice = append(valSlice, p.LoginName)
	}

	if p.LoginPass != "" {
		colNames += "login_pass,"
		colTags += "?,"
		valSlice = append(valSlice, p.LoginPass)
	}

	if p.Status != 0 {
		colNames += "status,"
		colTags += "?,"
		valSlice = append(valSlice, p.Status)
	}

	if p.AvatarUrl != "" {
		colNames += "avatar_url,"
		colTags += "?,"
		valSlice = append(valSlice, p.AvatarUrl)
	}

	if p.NickName != "" {
		colNames += "nick_name,"
		colTags += "?,"
		valSlice = append(valSlice, p.NickName)
	}

	if p.Gender != 0 {
		colNames += "gender,"
		colTags += "?,"
		valSlice = append(valSlice, p.Gender)
	}

	if p.City != "" {
		colNames += "city,"
		colTags += "?,"
		valSlice = append(valSlice, p.City)
	}

	if p.Province != "" {
		colNames += "province,"
		colTags += "?,"
		valSlice = append(valSlice, p.Province)
	}

	if p.Country != "" {
		colNames += "country,"
		colTags += "?,"
		valSlice = append(valSlice, p.Country)
	}

	if p.Language != "" {
		colNames += "language,"
		colTags += "?,"
		valSlice = append(valSlice, p.Language)
	}

	if p.Errors != 0 {
		colNames += "errors,"
		colTags += "?,"
		valSlice = append(valSlice, p.Errors)
	}

	if p.AccountBal != 0.00 {
		colNames += "account_bal,"
		colTags += "?,"
		valSlice = append(valSlice, p.AccountBal)
	}

	if p.Market != "" {
		colNames += "market,"
		colTags += "?,"
		valSlice = append(valSlice, p.Market)
	}

	if p.RandomNo != 0 {
		colNames += "random_no,"
		colTags += "?,"
		valSlice = append(valSlice, p.RandomNo)
	}

	if p.CreatedTime != "" {
		colNames += "created_time,"
		colTags += "?,"
		valSlice = append(valSlice, p.CreatedTime)
	}

	if p.UpdatedTime != "" {
		colNames += "updated_time,"
		colTags += "?,"
		valSlice = append(valSlice, p.UpdatedTime)
	}

	if p.Memo != "" {
		colNames += "memo,"
		colTags += "?,"
		valSlice = append(valSlice, p.Memo)
	}

	if p.Version != 0 {
		colNames += "version,"
		colTags += "?,"
		valSlice = append(valSlice, p.Version)
	}

	colNames = strings.TrimRight(colNames, ",")
	colTags = strings.TrimRight(colTags, ",")
	exeSql := fmt.Sprintf("Insert into  b_account(%s)  values(%s)", colNames, colTags)
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

func (r AccountList) InsertMap(m map[string]interface{}, tr *sql.Tx) error {
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

	exeSql := fmt.Sprintf("Insert into  b_account(%s)  values(%s)", colNames, colTags)
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

func (r AccountList) UpdataEntity(keyNo string, p Account, tr *sql.Tx) error {
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

	if p.WxOpenId != "" {
		colNames += "wx_open_id=?,"

		valSlice = append(valSlice, p.WxOpenId)
	}

	if p.WxUnionId != "" {
		colNames += "wx_union_id=?,"

		valSlice = append(valSlice, p.WxUnionId)
	}

	if p.WxSessionKey != "" {
		colNames += "wx_session_key=?,"

		valSlice = append(valSlice, p.WxSessionKey)
	}

	if p.LoginMode != 0 {
		colNames += "login_mode=?,"
		valSlice = append(valSlice, p.LoginMode)
	}

	if p.LoginName != "" {
		colNames += "login_name=?,"

		valSlice = append(valSlice, p.LoginName)
	}

	if p.LoginPass != "" {
		colNames += "login_pass=?,"

		valSlice = append(valSlice, p.LoginPass)
	}

	if p.Status != 0 {
		colNames += "status=?,"
		valSlice = append(valSlice, p.Status)
	}

	if p.AvatarUrl != "" {
		colNames += "avatar_url=?,"

		valSlice = append(valSlice, p.AvatarUrl)
	}

	if p.NickName != "" {
		colNames += "nick_name=?,"

		valSlice = append(valSlice, p.NickName)
	}

	if p.Gender != 0 {
		colNames += "gender=?,"
		valSlice = append(valSlice, p.Gender)
	}

	if p.City != "" {
		colNames += "city=?,"

		valSlice = append(valSlice, p.City)
	}

	if p.Province != "" {
		colNames += "province=?,"

		valSlice = append(valSlice, p.Province)
	}

	if p.Country != "" {
		colNames += "country=?,"

		valSlice = append(valSlice, p.Country)
	}

	if p.Language != "" {
		colNames += "language=?,"

		valSlice = append(valSlice, p.Language)
	}

	if p.Errors != 0 {
		colNames += "errors=?,"
		valSlice = append(valSlice, p.Errors)
	}

	if p.AccountBal != 0.00 {
		colNames += "account_bal=?,"
		valSlice = append(valSlice, p.AccountBal)
	}

	if p.Market != "" {
		colNames += "market=?,"

		valSlice = append(valSlice, p.Market)
	}

	if p.RandomNo != 0 {
		colNames += "random_no=?,"
		valSlice = append(valSlice, p.RandomNo)
	}

	if p.CreatedTime != "" {
		colNames += "created_time=?,"

		valSlice = append(valSlice, p.CreatedTime)
	}

	if p.UpdatedTime != "" {
		colNames += "updated_time=?,"

		valSlice = append(valSlice, p.UpdatedTime)
	}

	if p.Memo != "" {
		colNames += "memo=?,"

		valSlice = append(valSlice, p.Memo)
	}

	if p.Version != 0 {
		colNames += "version=?,"
		valSlice = append(valSlice, p.Version)
	}

	colNames = strings.TrimRight(colNames, ",")
	valSlice = append(valSlice, keyNo)

	exeSql := fmt.Sprintf("update  b_account  set %s  where id=? ", colNames)
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

func (r AccountList) UpdateMap(keyNo string, m map[string]interface{}, tr *sql.Tx) error {
	l := time.Now()

	var colNames string
	valSlice := make([]interface{}, 0)
	for k, v := range m {
		colNames += k + "=?,"
		valSlice = append(valSlice, v)
	}
	valSlice = append(valSlice, keyNo)
	colNames = strings.TrimRight(colNames, ",")
	updateSql := fmt.Sprintf("Update b_account set %s where id=?", colNames)
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

func (r AccountList) Delete(keyNo string, tr *sql.Tx) error {
	l := time.Now()
	delSql := fmt.Sprintf("Delete from  b_account  where id=?")
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
