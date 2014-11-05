package model

type TService struct {
	TType
	Extends   *TService
	Functions []*TFunction
}

func (s *TService) AddFunc(f *TFunction) {
	if s.Functions == nil {
		s.Functions = []*TFunction{f}
		return
	}
	s.Functions = append(s.Functions, f)
}
