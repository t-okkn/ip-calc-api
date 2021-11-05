package models

type QuestionSet struct {
	Id      string
	Number   int
	Souce    string
	CIDRbits int
	IsCIDR   bool
}

type AnswerSet struct {
	Id     string
	Number  int
	NwAddr  string
	BcAddr  string
	Elapsed int
}

