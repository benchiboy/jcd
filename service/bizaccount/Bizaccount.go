package bizaccount

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
	ContractId  string `json:"contract_id"`
	AcctNo      string `json:"acct_no"`
	SystemId    string `json:"system_id"`
	Pricipal    string `json:"delinquent_principal"`
	Fee         string `json:"delinquent_fine"`
	LoanCapital string `json:"loan_capital"`
	PayDate     string `json:"open_acct_date"`
	Penalty     string `json:"expense_amt"`
	Total       string `json:"last_cyc_stmt_bal"`
	OverdueDays string `json:"delinquent_days"`
	Interest    string `json:"stmt_interest_bal"`
	UpdateDate  string `json:"stra_upd_date"`
	DueDate     string `json:"statement_day"`

	PageNo     int    `json:"page_no"`
	PageSize   int    `json:"page_size"`
	ExtraWhere string `json:"extra_where"`
	SortFld    string `json:"sort_fld"`
}

type AccountList struct {
	DB      *sql.DB
	Level   int
	Total   int       `json:"total"`
	Account []Account `json:"Account"`
}

type Account struct {
	CrfUid      string `json:"crf_uid"`
	ContractId  string `json:"contract_id"`
	AcctNo      string `json:"acct_no"`
	SystemId    string `json:"system_id"`
	Pricipal    string `json:"delinquent_principal"`
	IntFeeAmt   string `json:"int_fee_amt"`
	LoanCapital string `json:"loan_capital"`
	PayDate     string `json:"open_acct_date"`
	Total       string `json:"last_cyc_stmt_bal"`
	OverdueDays string `json:"delinquent_days"`
	UpdateDate  string `json:"stra_upd_date"`
	DueDate     string `json:"statement_day"`
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
	return &AccountList{DB: db, Total: 0, Account: make([]Account, 0), Level: level}
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
	return &AccountList{DB: db, Total: 0, Account: make([]Account, 0), Level: level}
}

/*
	说明：得到符合条件的总条数
	入参：s: 查询条件
	出参：参数1：返回符合条件的总条件, 参数2：如果错误返回错误对象
*/

func (r *AccountList) GetTotal(s Search) (int, error) {
	var where string
	l := time.Now()

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}

	qrySql := fmt.Sprintf("Select count(1) as total from tbl_ccms_biz_cust_customer   where 1=1 %s", where)
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

//func (r AccountList) Get(s Search) (*Account, error) {
//	var where string
//	l := time.Now()

//	if s.ExtraWhere != "" {
//		where += s.ExtraWhere
//	}

//	if s.AcctNo != "" {
//		where += s.AcctNo
//	}

//	qrySql := fmt.Sprintf("Select crfuid,acct_no,system_id, delinquent_principal,delinquent_fine,loan_capital, open_acct_date,expense_amt,last_cyc_stmt_bal,delinquent_days,stmt_interest_bal,stra_upd_date,statement_day   from tbl_ccms_biz_cust_customer where 1=1 %s ", where)
//	if r.Level == DEBUG {
//		log.Println(SQL_SELECT, qrySql)
//	}
//	rows, err := r.DB.Query(qrySql)
//	if err != nil {
//		log.Println(SQL_ERROR, err.Error())
//		return nil, err
//	}
//	defer rows.Close()

//	var p Account
//	if !rows.Next() {
//		return nil, fmt.Errorf("Not Finded Record")
//	} else {
//		err := rows.Scan(&p.CrfUid, &p.ContractId, &p.SystemId, &p.Pricipal, &p.Fee, &p.LoanCapital, &p.PayDate, &p.Penalty, &p.Total, &p.OverdueDays, &p.Interest, &p.UpdateDate, &p.DueDate)
//		if err != nil {
//			log.Println(SQL_ERROR, err.Error())
//			return nil, err
//		}
//	}
//	log.Println(SQL_ELAPSED, r)
//	if r.Level == DEBUG {
//		log.Println(SQL_ELAPSED, time.Since(l))
//	}
//	return &p, nil
//}

/*
	说明：根据条件查询复核条件对象列表，支持分页查询
	入参：s: 查询条件
	出参：参数1：返回符合条件的对象列表, 参数2：如果错误返回错误对象
*/

func (r *AccountList) GetList(s Search) ([]Account, error) {
	var where string
	l := time.Now()

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}
	if s.CrfUid != "" {
		where += " and crfuid='" + fmt.Sprintf("%s'", s.CrfUid)
	}

	var qrySql string
	if s.PageSize == 0 && s.PageNo == 0 {
		qrySql = fmt.Sprintf("Select crfuid,acct_no,system_id, delinquent_principal,(delinquent_fine+expense_amt+stmt_interest_bal) as intfee,loan_capital, open_acct_date,last_cyc_stmt_bal,delinquent_days ,date(statement_day) from tbl_ccms_biz_acct_account where 1=1 %s", where)
	} else {
		qrySql = fmt.Sprintf("Select crfuid,acct_no,system_id, delinquent_principal,(delinquent_fine+expense_amt+stmt_interest_bal) as intfee ,loan_capital, open_acct_date,last_cyc_stmt_bal,delinquent_days,date(statement_day) from tbl_ccms_biz_acct_account where 1=1 %s Limit %d offset %d", where, s.PageSize, (s.PageNo-1)*s.PageSize)
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
		err := rows.Scan(&p.CrfUid, &p.ContractId, &p.SystemId, &p.Pricipal, &p.IntFeeAmt, &p.LoanCapital, &p.PayDate, &p.Total, &p.OverdueDays, &p.DueDate)
		fmt.Println(err)
		r.Account = append(r.Account, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.Account, nil
}
