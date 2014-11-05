package parse

type Document struct {
	Namespaces map[string]string
	Includes   []string
	Consts     []*Constant
}

type Constant struct {
	Name      string
	FieldType string
}

func NewDocument() *Document {
	return &Document{
		Namespaces: map[string]string{},
		Includes:   []string{},
	}
}

func (d *Document) AddConst(c *Constant) {
	d.Consts = append(d.Consts, c)
}
