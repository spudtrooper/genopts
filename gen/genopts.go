package gen

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spudtrooper/goutil/slice"
)

type reqField struct {
	name, typ string
}

type field struct {
	Name, Type, Default, DefaultSelector string
}

type function struct {
	FunctionName string
	Field        field
}

type interfaceFunction struct {
	ObjectName   string
	FunctionName string
	Field        field
	MaybeQuote   string // This is a hack to get around the fact that we need to quote strings in the or and json-default annotation.
}

type typeDef struct {
	name   string
	fields []field
}

func GenOpts(optType, implType string, dir, goImportsBin string, fieldDefs []string, optss ...GenOptsOption) (string, error) {
	opts := MakeGenOptsOptions(optss...)
	extends := removeQuotes(opts.Extends())

	dir = removeQuotes(dir)
	goImportsBin = removeQuotes(goImportsBin)
	fieldDefs = removeQuotesSlice(fieldDefs)

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
		s := strings.TrimSpace(s)
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

	var extendedTypes []typeDef
	if extends != "" {
		e, err := findExtendedTypes(dir, extends)
		if err != nil {
			return "", err
		}
		extendedTypes = e
	}

	var outfile string
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.Errorf("os.Getwd: %v", err)
	}
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

	abs, err := filepath.Abs(outfile)
	if err != nil {
		return "", errors.Errorf("filepath.Abs(%q): %v", outfile, err)
	}
	pkg := path.Base(path.Dir(abs))

	output, err := genOutput(pkg, optType, implType, fieldDefs, prefix, opts.GenerateParamsStruct(), reqFields, extendedTypes)
	if err != nil {
		return "", err
	}

	addCommandLine := !opts.Nocommandline() && opts.Function() == ""
	if err := outputResult(pkg, outfile, output, optType, originalImplType, addCommandLine, opts); err != nil {
		return "", err
	}
	if err := postGenCleanup(goImportsBin, dir, outfile); err != nil {
		return "", err
	}
	output = ""

	return output, nil
}

var (
	//go : generate genopts --function RestaurantDetails --params verbose debugFailure
	genOptsFnRE = regexp.MustCompile(`^//go.generate genopts (.*)`)
	fnExtRE     = regexp.MustCompile(`--function (\S+).*`)
)

func findExtendedTypes(dir string, extends string) ([]typeDef, error) {
	extendsNames := slice.Strings(extends, ",")

	var goFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if ext := filepath.Ext(f.Name()); ext == ".go" {
			goFiles = append(goFiles, f.Name())
		}
	}

	var cmdLinesToRun []string
	extendsMap := make(map[string]bool)
	for _, e := range extendsNames {
		extendsMap[e] = true
	}
	for _, f := range goFiles {
		c, err := ioutil.ReadFile(filepath.Join(dir, f))
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(string(c), "\n") {
			if m := genOptsFnRE.FindStringSubmatch(line); len(m) == 2 {
				cmdLine := m[1]
				var fn string
				if m := fnExtRE.FindStringSubmatch(cmdLine); len(m) == 2 {
					fn = m[1]
				}
				if _, ok := extendsMap[fn]; ok {
					cmdLinesToRun = append(cmdLinesToRun, cmdLine)
				}
			}

		}
	}

	var typeDefs []typeDef
	for _, cmdLine := range cmdLinesToRun {
		r := csv.NewReader(strings.NewReader(cmdLine))
		r.Comma = ' '
		args, err := r.Read()
		if err != nil {
			return nil, err
		}
		if err := exec.Command("genopts", args...).Run(); err != nil {
			return nil, err
		}

		rest, name := findRest(args)
		fields := makeFields(rest)
		t := typeDef{
			name:   name,
			fields: fields,
		}
		typeDefs = append(typeDefs, t)
	}

	return typeDefs, nil
}

func isArg(arg, name string) bool { return arg == "--"+name || arg == "-"+name }

func findRest(args []string) ([]string, string) {
	var res []string
	var name string
	for i := 0; i < len(args); {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			for j := i; j < len(args); j++ {
				res = append(res, args[j])
			}
			break
		}
		i++
		if isArg(arg, "function") {
			name = args[i]
			i++
		} else if isArg(arg, "opt_type") ||
			isArg(arg, "impl_type") ||
			isArg(arg, "prefix") ||
			isArg(arg, "function") ||
			isArg(arg, "outfile") ||
			isArg(arg, "update_dir") ||
			isArg(arg, "update_field") ||
			isArg(arg, "goimports") ||
			isArg(arg, "exclude_dirs") ||
			isArg(arg, "config") ||
			isArg(arg, "logfile") ||
			isArg(arg, "required") ||
			isArg(arg, "extends") {
			i++
		}
	}
	return res, name
}

func removeQuotesSlice(ss []string) []string {
	var res []string
	for _, s := range ss {
		res = append(res, removeQuotes(s))
	}
	return res
}

func outputResult(pkg, outfile, output, optType, implType string, addCommandLine bool, opts GenOptsOptions) error {
	const tmpl = `
// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package {{.Package}}

{{if .AddCommandLine}}
//go:` + `generate genopts {{.CommandLine}}
{{end}}

{{.Output}}
	`
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

func makeFields(fieldDefs []string) []field {
	var fields []field
	for _, f := range fieldDefs {
		parts := strings.Split(f, ":")
		name := parts[0]
		typ := "bool"
		if len(parts) >= 2 {
			name = parts[0]
			typ = parts[1]
		}
		f := field{
			Name: name,
			Type: typ,
		}
		if len(parts) >= 3 {
			defaultSelector := title(f.Type)
			switch f.Type {
			case "time.Time":
				defaultSelector = "Time"
			case "time.Duration":
				defaultSelector = "Duration"
			}
			f.DefaultSelector = defaultSelector
			f.Default = strings.Join(parts[2:], ":")
		}
		fields = append(fields, f)
	}
	return fields
}

func title(str string) string {
	if str == "" {
		return ""
	}
	s := []rune(str)
	s[0] = unicode.ToUpper(s[0])
	return string(s)
}

func genOutput(pkg, optType, implType string, fieldDefs []string, functionPrefix string, genParamsStruct bool, reqFields []reqField, extends []typeDef) (string, error) {
	const tmpl = `
{{$optType := .OptType}}
{{$package := .Package}}
{{$implType := .ImplType}}
{{$functionPrefix := .FunctionPrefix}}

import (
	"github.com/spudtrooper/goutil/or"
)

type {{.OptType}} struct {
	f func(*{{.ImplType}})
	s string
}

func (o {{.OptType}}) String() string { return o.s }

type {{.OptType}}s interface {
{{range .InterfaceFunctions}}	
{{.FunctionName}}() {{.Field.Type}}
Has{{.FunctionName}}() bool{{end}}
{{- range .ToTypes}}
	To{{.ReturnType}}s() []{{.ReturnType}}
{{- end}}
}

{{range .Functions}}
func {{.FunctionName}}({{.Field.Name}} {{.Field.Type}}) {{$optType}} {
	return {{$optType}}{func(opts *{{$implType}}) {
		opts.has_{{.Field.Name}} = true
		opts.{{.Field.Name}} = {{.Field.Name}}
	}, "{{$package}}.{{.FunctionName}}({{.Field.Type}})"}
}
func {{.FunctionName}}Flag({{.Field.Name}} *{{.Field.Type}}) {{$optType}} {
	return {{$optType}}{func(opts *{{$implType}}) {
		if {{.Field.Name}} == nil {
			return
		}
		opts.has_{{.Field.Name}} = true
		opts.{{.Field.Name}} = *{{.Field.Name}}
	}, "{{$package}}.{{.FunctionName}}({{.Field.Type}})"}
}
{{end}}
type {{.ImplType}} struct {
{{range .Fields}}	{{.Name}} {{.Type}}
has_{{.Name}} bool
{{end}}
}
{{- range .InterfaceFunctions}}
	{{- if .Field.Default}}
	{{- if eq .Field.Type "bool" }}
		func ({{.ObjectName}} *{{$implType}}) {{.FunctionName}}() {{.Field.Type}} {  
			if {{.ObjectName}}.Has{{.FunctionName}}() {
				return {{.ObjectName}}.{{.Field.Name}}
			}
			return {{.Field.Default}}
		}
	{{- else }}
		func ({{.ObjectName}} *{{$implType}}) {{.FunctionName}}() {{.Field.Type}} {  return or.{{.Field.DefaultSelector}}({{.ObjectName}}.{{.Field.Name}}, {{.MaybeQuote}}{{.Field.Default}}{{.MaybeQuote}}) }
	{{- end}}
	{{- else}}
	func ({{.ObjectName}} *{{$implType}}) {{.FunctionName}}() {{.Field.Type}} { return {{.ObjectName}}.{{.Field.Name}} }
	{{- end}}
	func ({{.ObjectName}} *{{$implType}}) Has{{.FunctionName}}() bool { return {{.ObjectName}}.has_{{.Field.Name}} }
{{- end}}

{{if .GenParamsStruct}}
	type  {{.ParamsStructName}} struct {
		{{range .RequiredFields}}{{.Name}} {{.Type}}` + " `" + `json:"{{.SnakeName}}"{{if .Required}} required:"true"{{end}}{{if .Default}} default:"{{.MaybeEscapedQuote}}{{.Default}}{{.MaybeEscapedQuote}}"{{end}}` + "`" + `
		{{end}}
	}

	func (o {{.ParamsStructName}}) Options() []{{.OptType}} {
		return []{{.OptType}}{
			{{- range .RequiredFields}}
				{{- if not .Required}}
					{{$functionPrefix}}{{.Name}}(o.{{.Name}}),
				{{- end}}
			{{- end}}
		}	
	}
{{end}}

{{range .ToTypes}}
	{{$prefix := .Prefix}}
	// To{{.ReturnType}}s converts {{$optType}} to an array of {{.ReturnType}}
	func (o *{{$implType}})To{{.ReturnType}}s() []{{.ReturnType}} {
		return []{{.ReturnType}} {
			{{- range .FieldNames}}
				{{$prefix}}{{.}}(o.{{.}}()),
			{{- end}}
		}
}
{{end}}

func make{{.ImplTypeCaps}}(opts ...{{.OptType}}) *{{.ImplType}} {
	res := &{{.ImplType}}{}
	for _, opt := range opts {
		opt.f(res)
	}
	return res
}

func Make{{.OptType}}s(opts ...{{.OptType}}) {{.OptType}}s {
	return make{{.ImplTypeCaps}}(opts...)
}
`

	snake := func(s string) string {
		return strcase.ToSnake(s)
	}

	fields := makeFields(fieldDefs)
	seenFields := map[string]bool{}
	for _, f := range fields {
		if seenFields[f.Name] {
			return "", fmt.Errorf("duplicate field %q", f.Name)
		}
		seenFields[f.Name] = true
	}

	for _, td := range extends {
		for _, f := range td.fields {
			if !seenFields[f.Name] {
				fields = append(fields, f)
				seenFields[f.Name] = true
			}
		}
	}

	type toType struct {
		Prefix     string
		ReturnType string
		FieldNames []string
	}
	var toTypes []toType
	for _, td := range extends {
		var fieldNames []string
		for _, f := range td.fields {
			fieldNames = append(fieldNames, title(f.Name))
		}
		tt := toType{
			Prefix:     td.name,
			ReturnType: td.name + "Option",
			FieldNames: fieldNames,
		}
		toTypes = append(toTypes, tt)
	}

	type requiredField struct {
		Name, Type, SnakeName string
		Required              bool
		Default               string
		MaybeEscapedQuote     string // This is a hack to get around the fact that we need to quote strings in the or and json-default annotation.
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
		rf := requiredField{
			Name:      title(f.Name),
			SnakeName: snake(f.Name),
			Type:      f.Type,
			Default:   f.Default,
		}
		if f.Type == "string" {
			rf.MaybeEscapedQuote = `\"`
		}
		requiredFields = append(requiredFields, rf)
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
		ifn := interfaceFunction{
			ObjectName:   strings.ToLower(string(implType[0])),
			FunctionName: title(f.Name),
			Field:        f,
		}
		if f.Type == "string" {
			ifn.MaybeQuote = `"`
		}
		interfaceFunctions = append(interfaceFunctions, ifn)

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
		FunctionPrefix     string
		ToTypes            []toType
		Package            string
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
		FunctionPrefix:     functionPrefix,
		ToTypes:            toTypes,
		Package:            pkg,
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
