package genopts

import (
	"bytes"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/spudtrooper/genopts/options"
)

func GenOpts(optsType, implType string, fieldDefs []string, opts ...options.Option) (string, error) {
	o := options.MakeOptions(opts...)
	if implType == "" {
		s := []rune(optsType)
		s[0] = unicode.ToLower(s[0])
		implType = string(s) + "Impl"
	}
	var prefix string
	if o.Prefix() != "" {
		prefix = o.Prefix()
	} else if o.PrefixOptsType() {
		prefix = optsType
	}
	output, err := genOpts(optsType, implType, fieldDefs, prefix)
	if err != nil {
		return "", err
	}
	return output, nil
}

func genOpts(optsType, implType string, fieldDefs []string, functionPrefix string) (string, error) {
	const tmpl = `
{{$optsType := .OptsType}}
{{$implType := .ImplType}}
type {{.OptsType}} func(*{{.ImplType}})

type {{.OptsType}}s interface {
{{range .InterfaceFunctions}}	{{.FunctionName}}() {{.Field.Type}}
{{end}}
}
{{range .Functions}}
func {{.FunctionName}}({{.Field.Name}} {{.Field.Type}}) {{$optsType}} {
	return func(opts *{{$implType}}) {
		opts.{{.Field.Name}} = {{.Field.Name}}
	}
}
{{end}}
type {{.ImplType}} struct {
{{range .Fields}}	{{.Name}} {{.Type}}
{{end}}
}
{{range .InterfaceFunctions}}
func ({{.ObjectName}} *{{$implType}}) {{.FunctionName}}() {{.Field.Type}} { return {{.ObjectName}}.{{.Field.Name}} }{{end}}

func make{{.ImplTypeCaps}}(opts ...{{.OptsType}}) {{.ImplType}} {
	var res {{.ImplType}}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`
	type field struct {
		Name, Type string
	}
	type function struct {
		FunctionName string
		Field        field
	}
	type interfaceFunction struct {
		ObjectName   string
		FunctionName string
		Field        field
	}

	title := func(str string) string {
		s := []rune(str)
		s[0] = unicode.ToUpper(s[0])
		return string(s)
	}

	var fields []field
	for _, f := range fieldDefs {
		parts := strings.Split(f, ":")
		name := parts[0]
		typ := "bool"
		if len(parts) == 2 {
			name = parts[0]
			typ = parts[1]
		}
		fields = append(fields, field{
			Name: name,
			Type: typ,
		})
	}

	var functions []function
	var interfaceFunctions []interfaceFunction
	for _, f := range fields {
		functionName := title(f.Name)
		if functionPrefix != "" {
			functionName = functionPrefix + functionName
		}
		functions = append(functions, function{
			FunctionName: functionName,
			Field:        f,
		})
		interfaceFunctions = append(interfaceFunctions, interfaceFunction{
			ObjectName:   strings.ToLower(string(implType[0])),
			FunctionName: title(f.Name),
			Field:        f,
		})
	}

	var buf bytes.Buffer
	if err := renderTemplate(&buf, tmpl, "tmpl", struct {
		OptsType           string
		ImplType           string
		ImplTypeCaps       string
		Functions          []function
		InterfaceFunctions []interfaceFunction
		Fields             []field
	}{
		OptsType:           optsType,
		ImplType:           implType,
		ImplTypeCaps:       title(implType),
		Functions:          functions,
		InterfaceFunctions: interfaceFunctions,
		Fields:             fields,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func renderTemplate(buf io.Writer, t string, name string, data interface{}) error {
	tmpl, err := template.New(name).Parse(strings.TrimSpace(t))
	if err != nil {
		return err
	}
	if err := tmpl.Execute(buf, data); err != nil {
		return err
	}
	return nil
}

func title(str string) string {
	s := []rune(str)
	s[0] = unicode.ToTitle(s[0])
	return string(s)
}
