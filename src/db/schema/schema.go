package schema

type Field struct {
	Name string
	Type string
}

type Index struct {
	Name string
}

type Schema struct {
	ClassName string
	Fields    []Field
	Indices   []Index
}
