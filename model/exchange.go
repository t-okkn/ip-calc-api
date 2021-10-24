package model

type QuestionSet struct {
	Key      string
	Number   int
	Souce    string
	CIDRbits int
	isCIDR   bool
}

type AnswerSet struct {
	Key    string
	Number int
	NwAddr string
	BcAddr string
}

