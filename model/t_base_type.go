package model

type TBase int

const (
	TYPE_VOID TBase = iota
	TYPE_STRING
	TYPE_BOOL
	TYPE_BYTE
	TYPE_I16
	TYPE_I32
	TYPE_I64
	TYPE_DOUBLE
)

type TBaseType struct {
	TType
	Base TBase
}
