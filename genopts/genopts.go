package genopts

import (
	"bytes"
	"io"
	"strings"
	"text/template"
	"unicode"
)

func GenOpts(optsType, implType string, fieldDefs []string) (string, error) {
	if implType == "" {
		s := []rune(optsType)
		s[0] = unicode.ToLower(s[0])
		implType = string(s) + "Impl"
	}
	output, err := genOpts(optsType, implType, fieldDefs)
	if err != nil {
		return "", err
	}
	return output, nil
}

func genOpts(optsType, implType string, fieldDefs []string) (string, error) {
	var buf bytes.Buffer

	const tmpl = `
{{$optsType := .OptsType}}
{{$implType := .ImplType}}
type {{.OptsType}} func(*{{.ImplType}})
{{range .Functions}}
func {{.FunctionName}}({{.Field}} {{.ParamType}}) {{$optsType}} {
	return func(opts *{{$implType}}) {
		opts.{{.Field}} = {{.Field}}
	}
}
{{end}}
type {{.ImplType}} struct {
{{range .Fields}}	{{.Name}} {{.Type}}
{{end}}
}
`
	type function struct {
		FunctionName string
		ParamType    string
		Field        string
	}
	type field struct {
		Name, Type string
	}
	var functions []function
	var fields []field
	for _, f := range fieldDefs {
		parts := strings.Split(f, ":")
		name := parts[0]
		typ := "bool"
		if len(parts) == 2 {
			name = parts[0]
			typ = parts[1]
		}
		functionName := title(name)
		functions = append(functions, function{
			FunctionName: functionName,
			Field:        name,
			ParamType:    typ,
		})
		fields = append(fields, field{
			Name: name,
			Type: typ,
		})
	}

	if err := renderTemplate(&buf, tmpl, "tmpl", struct {
		OptsType  string
		ImplType  string
		Functions []function
		Fields    []field
	}{
		OptsType:  optsType,
		ImplType:  implType,
		Functions: functions,
		Fields:    fields,
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
