package protosrc

type HttpRule struct {
	Method string
	Path   string
}

type Field struct {
	Name     string
	TypeName string
	Repeated bool
	MapKey   string
	MapValue string
}

type Message struct {
	Name   string
	Fields []Field
}

type EnumValue struct {
	Name  string
	Index int
}

type Enum struct {
	Name   string
	Values []EnumValue
}

type Rpc struct {
	Name     string
	Request  string
	Response string
	Http     HttpRule
}

type Service struct {
	Name string
	Rpcs []Rpc
}

type File struct {
	Path     string
	Package  string
	Imports  []string
	Messages []Message
	Enums    []Enum
	Services []Service
}

type Loaded struct {
	Root  *File
	Files map[string]*File
}

