package model

type TList struct {
	TType
	ElemType *TType
}

type TSet struct {
	TType
	ElemType *TType
}

type TMap struct {
	TType
	KeyType *TType
	ValType *TType
}
