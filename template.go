package main

import (
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

// Template suitable for executing
type Template interface {
	Execute(io.Writer, interface{}) error
}

// TemplateFromFile returns a new template created by parsing a file
func TemplateFromFile(filename string) (Template, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return TemplateFromReader(f)
}

// TemplateFromReader returns a new template created by parsing from a Reader
func TemplateFromReader(in io.Reader) (Template, error) {

	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	str := string(bytes)
	return TemplateFromString(str)
}

// TemplateFromString returns a new template created by parsing a string
func TemplateFromString(str string) (Template, error) {
	return newConcreteTemplate("str").Parse(str)
}

func newConcreteTemplate(tplName string) *template.Template {
	funcMap := template.FuncMap{"vault": vaultGetString}
	return template.New(tplName).Funcs(funcMap)
}
