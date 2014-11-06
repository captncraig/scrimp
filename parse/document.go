package parse

type Document struct {
	Namespaces map[string]string
	Includes   []string
	Consts     []*Constant
	Typedefs   []*Typedef
}

type Constant struct {
	Name      string
	FieldType string
	Value     string
}

type Typedef struct {
	Name      string
	FieldType string
}

func NewDocument() *Document {
	return &Document{
		Namespaces: map[string]string{},
		Includes:   []string{},
		Typedefs:   []*Typedef{},
		Consts:     []*Constant{},
	}
}

func (d *Document) AddConst(c *Constant) {
	d.Consts = append(d.Consts, c)
}

func (d *Document) AddTypedef(t *Typedef) {
	d.Typedefs = append(d.Typedefs, t)
}
