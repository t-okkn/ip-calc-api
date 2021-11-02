package models

type QuestionSet struct {
	Key      string
	Number   int
	Souce    string
	CIDRbits int
	IsCIDR   bool
}

type AnswerSet struct {
	Key     string
	Number  int
	NwAddr  string
	BcAddr  string
	Elapsed int
}

