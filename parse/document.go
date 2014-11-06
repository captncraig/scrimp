package parse

type Document struct {
	Namespaces map[string]string
	Includes   []string
	Consts     []*Constant
	Typedefs   []*Typedef
	Enums      []*Enum
	Structs    []*Struct
	Xceptions  []*Struct
	Services   []*Service
	DocText    string
}

type Constant struct {
	Name      string
	FieldType string
	Value     string
	DocText   string
}

type Typedef struct {
	Name      string
	FieldType string
	DocText   string
}

type Enum struct {
	Name    string
	Members map[string]int
	DocText string
}

type Struct struct {
	Name    string
	Fields  []*Field
	DocText string
}

func (s *Struct) AddField(f *Field) {
	if s.Fields == nil {
		s.Fields = []*Field{f}
		return
	}
	s.Fields = append(s.Fields, f)
}

type Field struct {
	Index        int
	Name         string
	FieldType    string
	DefaultValue string
	Required     bool
	DocText      string
}

type Service struct {
	Name      string
	Extends   string
	Functions []*Function
	DocText   string
}

type Function struct {
	Oneway     bool
	ReturnType string
	Name       string
	Fields     []*Field
	Throws     []*Field
	DocText    string
}

func (s *Service) AddFunction(f *Function) {
	s.Functions = append(s.Functions, f)
}

func NewDocument() *Document {
	return &Document{
		Namespaces: map[string]string{},
		Includes:   []string{},
		Typedefs:   []*Typedef{},
		Consts:     []*Constant{},
		Enums:      []*Enum{},
		Structs:    []*Struct{},
		Xceptions:  []*Struct{},
		Services:   []*Service{},
	}
}

func (d *Document) AddConst(c *Constant) {
	d.Consts = append(d.Consts, c)
}

func (d *Document) AddStruct(c *Struct) {
	d.Structs = append(d.Structs, c)
}

func (d *Document) AddXception(c *Struct) {
	d.Xceptions = append(d.Xceptions, c)
}

func (d *Document) AddEnum(e *Enum) {
	d.Enums = append(d.Enums, e)
}

func (d *Document) AddTypedef(t *Typedef) {
	d.Typedefs = append(d.Typedefs, t)
}

func (d *Document) AddService(s *Service) {
	d.Services = append(d.Services, s)
}
