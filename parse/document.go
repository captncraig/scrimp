package parse

type Document struct {
	Namespaces map[string]string
	Includes   []string
}

func NewDocument() *Document {
	return &Document{
		Namespaces: map[string]string{},
		Includes:   []string{},
	}
}
