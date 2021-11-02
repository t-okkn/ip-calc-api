package models

import (
	"github.com/go-gorp/gorp"
)

type MstrKey struct {
	Key    string `db:"key"`
	Total  int    `db:"total"`
	Expire string `db:"expire"`
}

type TranQuestion struct {
	Key       string `db:"key"`
	Number    int    `db:"question_number"`
	Source    string `db:"source"`
	CIDRbits  int    `db:"cidr_bits"`
	IsCIDR    bool   `db:"is_cidr"`
	CorNwAddr string `db:"correct_nw"`
	AnsNwAddr string `db:"answer_nw"`
	CorBcAddr string `db:"correct_bc"`
	AnsBcAddr string `db:"answer_bc"`
	Elapsed   int    `db:"elapsed"`
}

// MapStructsToTables 構造体と物理テーブルの紐付け
func MapStructsToTables(dbmap *gorp.DbMap) {
	dbmap.AddTableWithName(MstrKey{}, "M_KEY")
	dbmap.AddTableWithName(TranQuestion{}, "T_QUESTION")
}
