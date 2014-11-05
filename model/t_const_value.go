package model

type TConstValueType int

const (
	CV_INTEGER TConstValueType = iota
	CV_DOUBLE
	CV_STRING
	CV_MAP
	CV_LIST
	CV_IDENTIFIER
)

type TConstValue struct {
	MapVal    map[*TConstValue]*TConstValue
	ListVal   []*TConstValue
	IntVal    int64
	DoubleVal float64
	StringVal string
	Type      TConstValueType
	Enum      TEnum
}
