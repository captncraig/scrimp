package model

type TFunction struct {
	TDoc
	Return    *TType
	Name      string
	ArgList   *TStruct
	Xceptions *TStruct
	OneWay    bool
}
