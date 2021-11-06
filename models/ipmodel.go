package models

import (
	"github.com/go-gorp/gorp"
)

type MstrID struct {
	Id     string `db:"id, primarykey"`
	Total  int    `db:"total"`
	Expire string `db:"expire"`
}

type TranQuestion struct {
	Id        string `db:"id, primarykey"`
	Number    int    `db:"question_number, primarykey"`
	Source    string `db:"source"`
	CIDRbits  int    `db:"cidr_bits"`
	IsCIDR    int    `db:"is_cidr"`
	CorNwAddr string `db:"correct_nw"`
	AnsNwAddr string `db:"answer_nw"`
	CorBcAddr string `db:"correct_bc"`
	AnsBcAddr string `db:"answer_bc"`
	Elapsed   int    `db:"elapsed"`
}

// MapStructsToTables 構造体と物理テーブルの紐付け
func MapStructsToTables(dbmap *gorp.DbMap) {
	dbmap.AddTableWithName(MstrID{}, "M_ID").SetKeys(false, "Id")
	dbmap.AddTableWithName(TranQuestion{}, "T_QUESTION").SetKeys(false, "Id", "Number")
}

