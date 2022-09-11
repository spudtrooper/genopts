package gen

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/pkg/errors"
)

func GenOpts(optType, implType string, dir, goImportsBin string, fieldDefs []string, opts ...GenOptsOption) (string, error) {
	dir = removeQuotes(dir)
	goImportsBin = removeQuotes(goImportsBin)
	fieldDefs = removeQuotesSlice(fieldDefs)

	o := MakeGenOptsOptions(opts...)
	originalImplType := implType
	var prefix string
	if o.Prefix() != "" {
		prefix = o.Prefix()
		optType = prefix + "Option"
	} else if o.Function() != "" {
		prefix = o.Function()
		optType = prefix + "Option"
	} else if o.PrefixOptsType() {
		prefix = optType
	}
	if optType == "" {
		optType = "Option"
	}
	if implType == "" {
		s := []rune(optType)
		s[0] = unicode.ToLower(s[0])
		implType = string(s) + "Impl"
	}
	output, err := genOpts(optType, implType, fieldDefs, prefix)
	if err != nil {
		return "", err
	}

	var outfile string
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.Errorf("os.Getwd: %v", err)
	}
	log.Printf("have pwd: %s", pwd)
	if o.Outfile() != "" {
		// If the dirname of the outfile ends with the end of the pwd, then we are running in go generate mode
		// In this case, we use the basename of the outfile.
		if tailPwd, startOutfile := path.Base(pwd), path.Base(path.Dir(o.Outfile())); tailPwd == startOutfile {
			outfile = path.Base(o.Outfile())
		} else {
			outfile = o.Outfile()
		}
	} else {
		filename := strings.ToLower(prefix + "options.go")
		outfile = path.Join(pwd, filename)
	}

	addCommandLine := !o.Nocommandline() && o.Function() == ""
	if err := outputResult(outfile, output, optType, originalImplType, addCommandLine, o); err != nil {
		return "", err
	}
	if err := postGenCleanup(goImportsBin, dir, outfile); err != nil {
		return "", err
	}
	output = ""

	return output, nil
}

func removeQuotesSlice(ss []string) []string {
	var res []string
	for _, s := range ss {
		res = append(res, removeQuotes(s))
	}
	return res
}

func outputResult(outfile, output, optType, implType string, addCommandLine bool, opts GenOptsOptions) error {
	const tmpl = `
// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package {{.Package}}

{{if .AddCommandLine}}
//go:` + `generate genopts {{.CommandLine}}
{{end}}

{{.Output}}
	`

	abs, err := filepath.Abs(outfile)
	if err != nil {
		return errors.Errorf("filepath.Abs(%q): %v", outfile, err)
	}
	pkg := path.Base(path.Dir(abs))
	var cmdLineParts []string
	// This has to stay in sync with flags
	if optType != "Option" && optType != opts.Prefix()+"Option" { // The defaults
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
		cmdLineParts = append(cmdLineParts, fmt.Sprintf("\"%s\"", removeQuotes(fs)))
	}
	cmdLine := strings.Join(cmdLineParts, " ")

	var buf bytes.Buffer
	if err := renderTemplate(&buf, tmpl, "outputResult", struct {
		Package        string
		CommandLine    string
		Output         string
		AddCommandLine bool
	}{
		Package:        pkg,
		CommandLine:    cmdLine,
		Output:         output,
		AddCommandLine: addCommandLine,
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
{{range .InterfaceFunctions}}	
{{.FunctionName}}() {{.Field.Type}}
Has{{.FunctionName}}() bool{{end}}
}
{{range .Functions}}
func {{.FunctionName}}({{.Field.Name}} {{.Field.Type}}) {{$optType}} {
	return func(opts *{{$implType}}) {
		opts.has_{{.Field.Name}} = true
		opts.{{.Field.Name}} = {{.Field.Name}}
	}
}
func {{.FunctionName}}Flag({{.Field.Name}} *{{.Field.Type}}) {{$optType}} {
	return func(opts *{{$implType}}) {
		if {{.Field.Name}} == nil {
			return
		}
		opts.has_{{.Field.Name}} = true
		opts.{{.Field.Name}} = *{{.Field.Name}}
	}
}
{{end}}
type {{.ImplType}} struct {
{{range .Fields}}	{{.Name}} {{.Type}}
has_{{.Name}} bool
{{end}}
}
{{range .InterfaceFunctions}}
func ({{.ObjectName}} *{{$implType}}) {{.FunctionName}}() {{.Field.Type}} { return {{.ObjectName}}.{{.Field.Name}} }
func ({{.ObjectName}} *{{$implType}}) Has{{.FunctionName}}() bool { return {{.ObjectName}}.has_{{.Field.Name}} }{{end}}

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
