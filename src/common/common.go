package common

type Field struct {
	Name string
	Tags map[string]string
	Values map[string]string // we just assume they are strings, this is just backup app
}
