package models

type QuestionSet struct {
	Id         string `json:"id"`
	Number     int    `json:"number"`
	Source     string `json:"source"`
	CIDRbits   int    `json:"cidr_bits"`
	IsCIDR     int    `json:"is_cidr"`
	SubnetMask string `json:"subnet_mask"`
}

type AnswerSet struct {
	Id      string `json:"id"`
	Number  int    `json:"number"`
	NwAddr  string `json:"answer_nw"`
	BcAddr  string `json:"answer_bc"`
	Elapsed int    `json:"elapsed"`
}

