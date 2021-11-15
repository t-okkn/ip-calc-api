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
	NwAddr1st  string `form:"nwaddr_1st"`
	NwAddr2nd  string `form:"nwaddr_2nd"`
	NwAddr3rd  string `form:"nwaddr_3rd"`
	NwAddr4th  string `form:"nwaddr_4th"`
	BcAddr1st  string `form:"bcaddr_1st"`
	BcAddr2nd  string `form:"bcaddr_2nd"`
	BcAddr3rd  string `form:"bcaddr_3rd"`
	BcAddr4th  string `form:"bcaddr_4th"`
	Elapsed    string `form:"elapsed"`
}

