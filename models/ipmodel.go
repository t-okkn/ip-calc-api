package models

import (
	"github.com/go-gorp/gorp"
)

type MstrID struct {
	Id     string `db:"id, primarykey" json:"id"`
	Total  int    `db:"total" json:"total"`
	Expire string `db:"expire" json:"expire"`
}

type TranQuestion struct {
	Id        string `db:"id, primarykey" json:"id"`
	Number    int    `db:"question_number, primarykey" json:"number"`
	Source    string `db:"source" json:"source"`
	CIDRbits  int    `db:"cidr_bits" json:"cidr_bits"`
	IsCIDR    int    `db:"is_cidr" json:"is_cidr"`
	CorNwAddr string `db:"correct_nw" json:"correct_nw"`
	AnsNwAddr string `db:"answer_nw" json:"answer_nw"`
	CorBcAddr string `db:"correct_bc" json:"correct_bc"`
	AnsBcAddr string `db:"answer_bc" json:"answer_bc"`
	Elapsed   int    `db:"elapsed" json:"elapsed"`
}

// MapStructsToTables 構造体と物理テーブルの紐付け
func MapStructsToTables(dbmap *gorp.DbMap) {
	dbmap.AddTableWithName(MstrID{}, "M_ID").SetKeys(false, "Id")
	dbmap.AddTableWithName(TranQuestion{}, "T_QUESTION").SetKeys(false, "Id", "Number")
}

