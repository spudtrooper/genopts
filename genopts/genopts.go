package genopts

import (
	"bytes"
	"io"
	"strings"
	"text/template"
	"unicode"
)

// START-PASTE
type Option func(*optionImpl)

func PrefixOptsType(prefixOptsType bool) Option {
	return func(opts *optionImpl) {
		opts.prefixOptsType = prefixOptsType
	}
}

func Prefix(prefix string) Option {
	return func(opts *optionImpl) {
		opts.prefix = prefix
	}
}

type optionImpl struct {
	prefixOptsType bool
	prefix         string
}

func makeOptionImpl(opts ...Option) optionImpl {
	var res optionImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}

// END-PASTE

func GenOpts(optsType, implType string, fieldDefs []string, opts ...Option) (string, error) {
	o := makeOptionImpl(opts...)
	if implType == "" {
		s := []rune(optsType)
		s[0] = unicode.ToLower(s[0])
		implType = string(s) + "Impl"
	}
	var prefix string
	if o.prefix != "" {
		prefix = o.prefix
	} else if o.prefixOptsType {
		prefix = optsType
	}
	output, err := genOpts(optsType, implType, fieldDefs, prefix)
	if err != nil {
		return "", err
	}
	return output, nil
}

func genOpts(optsType, implType string, fieldDefs []string, functionPrefix string) (string, error) {
	var buf bytes.Buffer

	const tmpl = `
{{$optsType := .OptsType}}
{{$implType := .ImplType}}
type {{.OptsType}} func(*{{.ImplType}})
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
	for _, f := range fields {
		functionName := title(f.Name)
		if functionPrefix != "" {
			functionName = functionPrefix + functionName
		}
		functions = append(functions, function{
			FunctionName: functionName,
			Field:        f,
		})

	}

	var implTypeCaps string
	{
		s := []rune(implType)
		s[0] = unicode.ToUpper(s[0])
		implTypeCaps = string(s)
	}
	if err := renderTemplate(&buf, tmpl, "tmpl", struct {
		OptsType     string
		ImplType     string
		ImplTypeCaps string
		Functions    []function
		Fields       []field
	}{
		OptsType:     optsType,
		ImplType:     implType,
		ImplTypeCaps: implTypeCaps,
		Functions:    functions,
		Fields:       fields,
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
