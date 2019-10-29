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
	LoanDays    string `json:"loan_days"`
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

func (r *AccountList) GetList(s Search) ([]Account, error) {
	var where string
	l := time.Now()

	if s.ExtraWhere != "" {
		where += s.ExtraWhere
	}
	if s.CrfUid != "" {
		where += " and crf_uid='" + fmt.Sprintf("%s'", s.CrfUid)
	}

	var qrySql string
	qrySql = fmt.Sprintf("SELECT crf_uid,loan_no,mct_no,principal_amt,fee_amt+int_amt+penalty_amt as intfee,pay_amt,pay_date,principal_amt+int_amt+fee_amt+penalty_amt AS total_amt ,datediff(now(),curr_bill_date) as overdue_days,curr_bill_date,loan_periods FROM act_mv_loan_info  where 1=1 %s", where)

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
		err := rows.Scan(&p.CrfUid, &p.ContractId, &p.SystemId, &p.Pricipal, &p.IntFeeAmt, &p.LoanCapital, &p.PayDate, &p.Total, &p.OverdueDays, &p.DueDate, &p.LoanDays)
		fmt.Println(err)
		r.Account = append(r.Account, p)
	}
	log.Println(SQL_ELAPSED, r)
	if r.Level == DEBUG {
		log.Println(SQL_ELAPSED, time.Since(l))
	}
	return r.Account, nil
}
