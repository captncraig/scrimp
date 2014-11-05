package model

type TDoc struct {
	HasDoc bool
	Doc    string
}

func (t *TDoc) SetDoc(d string) {
	t.HasDoc = true
	t.Doc = d
}
