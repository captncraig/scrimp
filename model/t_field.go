package model

type EReq int

const (
	T_REQUIRED EReq = iota
	T_OPTIONAL
	T_OPT_IN_REQ_OUT
)

type TField struct {
	TDoc
	Type  *TType
	Name  string
	Key   int
	Req   EReq
	Value TConstValue
}
