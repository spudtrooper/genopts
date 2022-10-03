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

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

type reqField struct {
	name, typ string
}

func GenOpts(optType, implType string, dir, goImportsBin string, fieldDefs []string, optss ...GenOptsOption) (string, error) {
	dir = removeQuotes(dir)
	goImportsBin = removeQuotes(goImportsBin)
	fieldDefs = removeQuotesSlice(fieldDefs)

	opts := MakeGenOptsOptions(optss...)
	originalImplType := implType
	var prefix string
	if opts.Prefix() != "" {
		prefix = opts.Prefix()
		optType = prefix + "Option"
	} else if opts.Function() != "" {
		prefix = opts.Function()
		optType = prefix + "Option"
	} else if opts.PrefixOptsType() {
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

	var reqFields []reqField
	for _, s := range strings.Split(opts.RequiredFields(), ",") {
		if s == "" {
			continue
		}
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			parts = strings.Split(s, " ")
		}
		if len(parts) != 2 {
			return "", errors.Errorf("invalid required field %v in %q",
				parts, opts.RequiredFields())
		}
		reqFields = append(reqFields, reqField{
			name: parts[0],
			typ:  parts[1],
		})
	}

	output, err := genOutput(
		optType, implType, fieldDefs, prefix, opts.GenerateParamsStruct(), reqFields)
	if err != nil {
		return "", err
	}

	var outfile string
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.Errorf("os.Getwd: %v", err)
	}
	log.Printf("have pwd: %s", pwd)
	if opts.Outfile() != "" {
		// If the dirname of the outfile ends with the end of the pwd, then we are running in go generate mode
		// In this case, we use the basename of the outfile.
		if tailPwd, startOutfile := path.Base(pwd), path.Base(path.Dir(opts.Outfile())); tailPwd == startOutfile {
			outfile = path.Base(opts.Outfile())
		} else {
			outfile = opts.Outfile()
		}
	} else {
		filename := strings.ToLower(prefix + "options.go")
		outfile = path.Join(pwd, filename)
	}

	addCommandLine := !opts.Nocommandline() && opts.Function() == ""
	if err := outputResult(outfile, output, optType, originalImplType, addCommandLine, opts); err != nil {
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
	if opts.HasRequiredFields() && opts.RequiredFields() != "" {
		cmdLineParts = append(cmdLineParts, "--required", fmt.Sprintf("%q", opts.RequiredFields()))
	}
	if opts.GenerateParamsStruct() {
		cmdLineParts = append(cmdLineParts, "--params")
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

func genOutput(optType, implType string, fieldDefs []string, functionPrefix string, genParamsStruct bool, reqFields []reqField) (string, error) {
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

{{if .GenParamsStruct}}
type  {{.ParamsStructName}} struct {
	{{range .RequiredFields}}{{.Name}} {{.Type}}{{if .Required}}` + "`" + `json:"{{.SnakeName}}" required:"true"` + "`" + `{{else}}` + "`" + `json:"{{.SnakeName}}"` + "`" + `{{end}}
	{{end}}
}
{{end}}

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

	snake := func(s string) string {
		return strcase.ToSnake(s)
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

	type requiredField struct {
		Name, Type, SnakeName string
		Required              bool
	}
	var requiredFields []requiredField
	for _, f := range reqFields {
		requiredFields = append(requiredFields, requiredField{
			Name:      title(f.name),
			SnakeName: snake(f.name),
			Type:      f.typ,
			Required:  true,
		})
	}
	for _, f := range fields {
		requiredFields = append(requiredFields, requiredField{
			Name:      title(f.Name),
			SnakeName: snake(f.Name),
			Type:      f.Type,
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

	paramsStructName := title(strings.Replace(optType, "Option", "", 1)) + "Params"

	var buf bytes.Buffer
	if err := renderTemplate(&buf, tmpl, "genOpts", struct {
		OptType            string
		ImplType           string
		ImplTypeCaps       string
		Functions          []function
		InterfaceFunctions []interfaceFunction
		Fields             []field
		GenParamsStruct    bool
		ParamsStructName   string
		RequiredFields     []requiredField
	}{
		OptType:            optType,
		ImplType:           implType,
		ImplTypeCaps:       title(implType),
		Functions:          functions,
		InterfaceFunctions: interfaceFunctions,
		Fields:             fields,
		GenParamsStruct:    genParamsStruct,
		ParamsStructName:   paramsStructName,
		RequiredFields:     requiredFields,
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
