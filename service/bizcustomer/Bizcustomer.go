package bizcustomer

import (
	"database/sql"
	"fmt"
	"log"
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
	CrfUid      string `json:"crf_uid"`
	CertNo      string `json:"cert_no"`
	CustName    string `json:"cust_name"`
	MCardMobile string `json:"mcard_mobile"`

	PageNo     int    `json:"page_no"`
	PageSize   int    `json:"page_size"`
	ExtraWhere string `json:"extra_where"`
	SortFld    string `json:"sort_fld"`
}

type CustomerList struct {
	DB        *sql.DB
	Level     int
	Total     int        `json:"total"`
	Customers []Customer `json:"Customer"`
}

type Customer struct {
	CrfUid             string `json:"crf_uid"`
	ContractId         string `json:"contract_id"`
	CertNo             string `json:"cert_no"`
	CustName           string `json:"cust_name"`
	MCardMobile        string `json:"mcard_mobile"`
	BankCardNo         string `json:"bank_card_no"`
	MainContactAddr    string `json:"main_contact_addr"`
	UrgentLinkManName  string `json:"urgent_linkman_name"`
	UrgentLinkManPhone string `json:"urgent_linkman_phone"`
}

type Form struct {
	Form Customer `json:"Customer"`
}

/*
	说明：创建实例对象
	入参：db:数据库sql.DB, 数据库已经连接, level:日志级别
	出参：实例对象
*/

func New(db *sql.DB, level int) *CustomerList {
	if db == nil {
		log.Println(SQL_SELECT, "Database is nil")
		return nil
	}
	return &CustomerList{DB: db, Total: 0, Customers: make([]Customer, 0), Level: level}
}

/*
	说明：创建实例对象
	入参：url:连接数据的url, 数据库还没有CONNECTED, level:日志级别
	出参：实例对象
*/

func NewUrl(url string, level int) *CustomerList {
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
	return &CustomerList{DB: db, Total: 0, Customers: make([]Customer, 0), Level: level}
}

/*
	说明：得到符合条件的总条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的总条件, 参数2：如果错误返回错误对象
*/

func (r *CustomerList) GetTotal(s Search) (int, error) {
	var where string
	l := time.Now()

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	if s.CustName != "" {
		where += s.CustName
	}

	if s.CertNo != "" {
		where += s.CertNo
	}

	if s.MCardMobile != "" {
		where += s.MCardMobile
	}

	qrySql := fmt.Sprintf("Select count(1) as total from act_mv_borrower   where 1=1 %s", where)
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

func (r CustomerList) Get(s Search) (*Customer, error) {
	var where string
	l := time.Now()

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	if s.CustName != "" {
		where += s.CustName
	}

	if s.CertNo != "" {
		where += s.CertNo
	}

	if s.MCardMobile != "" {
		where += s.MCardMobile
	}

	qrySql := fmt.Sprintf("Select crfuid,,user_id_no,full_name,user_mobile,user_card_no  from act_mv_borrower where 1=1 %s ", where)

	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p Customer
	if !rows.Next() {
		return nil, fmt.Errorf("Not Finded Record")
	} else {
		err := rows.Scan(&p.CrfUid, &p.CertNo, &p.CustName, &p.MCardMobile, &p.BankCardNo)
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

func (r *CustomerList) GetList(s Search) ([]Customer, error) {
	var where string
	l := time.Now()

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}
	if s.CustName != "" {
		where += " and full_name='" + fmt.Sprintf("%s'", s.CustName)
	}

	if s.CertNo != "" {
		where += " and user_id_no='" + fmt.Sprintf("%s'", s.CertNo)
	}

	if s.MCardMobile != "" {
		where += " and user_mobile='" + fmt.Sprintf("%s'", s.MCardMobile)
	}

	qrySql := fmt.Sprintf("Select crf_uid,user_id_no,full_name,user_mobile,user_card_no  from act_mv_borrower where 1=1 %s ", where)

	if r.Level == DEBUG {
		log.Println(SQL_SELECT, qrySql)
	}
	rows, err := r.DB.Query(qrySql)
	if err != nil {
		log.Println(SQL_ERROR, err.Error())
		return nil, err
	}
	defer rows.Close()

	var p Customer
	for rows.Next() {
		rows.Scan(&p.CrfUid, &p.CertNo, &p.CustName, &p.MCardMobile, &p.BankCardNo)
		r.Customers = append(r.Customers, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.Customers, nil
}
