package model

type KeyMgmt struct {
	Key    string
	Total  int
	Expire string
}

type QuestionMgmt struct {
	Key       string
	Number    int
	Source    string
	CIDRbits  int
	isCIDR    bool
	CorNwAddr string
	AnsNwAddr string
	CorBcAddr string
	AnsBcAddr string
	Elapsed   int
}

