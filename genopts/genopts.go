package genopts

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"text/template"
	"unicode"

	"github.com/spudtrooper/genopts/options"
)

func GenOpts(optType, implType string, fieldDefs []string, opts ...options.Option) (string, error) {
	o := options.MakeOptions(opts...)
	originalImplType := implType
	if implType == "" {
		s := []rune(optType)
		s[0] = unicode.ToLower(s[0])
		implType = string(s) + "Impl"
	}
	var prefix string
	if o.Prefix() != "" {
		prefix = o.Prefix()
	} else if o.PrefixOptsType() {
		prefix = optType
	}
	output, err := genOpts(optType, implType, fieldDefs, prefix)
	if err != nil {
		return "", err
	}
	if o.Outfile() != "" {
		if err := outputResult(o.Outfile(), output, optType, originalImplType, o); err != nil {
			return "", err
		}
		output = ""
	}

	return output, nil
}

func outputResult(outfile, output, optType, implType string, opts options.Options) error {
	const tmpl = `
package {{.Package}}

// genopts {{.CommandLine}}

{{.Output}}
	`

	pkg := path.Base(path.Dir(outfile))
	var cmdLineParts []string
	// This has to stay in sync with flags
	if optType != "Option" { // The default
		cmdLineParts = append(cmdLineParts, "--opt_type="+optType)
	}
	if implType != "" {
		cmdLineParts = append(cmdLineParts, "--impl_type="+implType)
	}
	if opts.Prefix() != "" {
		cmdLineParts = append(cmdLineParts, "--prefix="+opts.Prefix())
	}
	if opts.PrefixOptsType() {
		cmdLineParts = append(cmdLineParts, "--prefix_opts_type")
	}
	cmdLineParts = append(cmdLineParts, "--outfile="+outfile)
	for _, fs := range flag.CommandLine.Args() {
		cmdLineParts = append(cmdLineParts, "'"+fs+"'")
	}
	cmdLine := strings.Join(cmdLineParts, " ")

	var buf bytes.Buffer
	if err := renderTemplate(&buf, tmpl, "outputResult", struct {
		Package     string
		CommandLine string
		Output      string
	}{
		Package:     pkg,
		CommandLine: cmdLine,
		Output:      output,
	}); err != nil {
		return err
	}

	if err := ioutil.WriteFile(outfile, buf.Bytes(), 0755); err != nil {
		return err
	}
	log.Printf("wrote to %s", outfile)

	return nil
}

func genOpts(optType, implType string, fieldDefs []string, functionPrefix string) (string, error) {
	const tmpl = `
{{$optType := .OptType}}
{{$implType := .ImplType}}
type {{.OptType}} func(*{{.ImplType}})

type {{.OptType}}s interface {
{{range .InterfaceFunctions}}	{{.FunctionName}}() {{.Field.Type}}
{{end}}
}
{{range .Functions}}
func {{.FunctionName}}({{.Field.Name}} {{.Field.Type}}) {{$optType}} {
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

func make{{.ImplTypeCaps}}(opts ...{{.OptType}}) *{{.ImplType}} {
	res := &{{.ImplType}}{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func Make{{.OptType}}s(opts ...{{.OptType}}) {{.OptType}}s {
	return make{{.ImplTypeCaps}}(opts...)
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
	if err := renderTemplate(&buf, tmpl, "genOpts", struct {
		OptType            string
		ImplType           string
		ImplTypeCaps       string
		Functions          []function
		InterfaceFunctions []interfaceFunction
		Fields             []field
	}{
		OptType:            optType,
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
