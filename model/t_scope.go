package model

type TScope struct {
	Types    map[string]*TType
	Consts   map[string]*TConst
	Services map[string]*TService
}
