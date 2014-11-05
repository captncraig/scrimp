package model

type TProgram struct {
	TDoc
	Path            string
	Name            string
	OutPath         string
	OutPathAbsolute bool
	Namespace       string
	Includes        []*TProgram
	IncludePrefix   string
	Typedefs        []*TTypedef
	Enums           []*TEnum
	Consts          []*TConst
	Objects         []*TStruct
	Structs         []*TStruct
	Xceptions       []*TStruct
	Services        []*TService
	Namespaces      map[string]string
}

func (t *TProgram) AddService(s *TService) {
	if t.Services == nil {
		t.Services = []*TService{s}
		return
	}
	t.Services = append(t.Services, s)
}
