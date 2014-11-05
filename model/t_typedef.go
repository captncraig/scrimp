package model

type TTypedef struct {
	Type     *TType
	Symbolic string
	Forward  bool
	Seen     bool
}
