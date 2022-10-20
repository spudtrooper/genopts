// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package gen

import "fmt"

//go:generate genopts --prefix=GenOpts --outfile=genoptsoptions.go "prefixOptsType:bool" "prefix:string" "function:string" "outfile:string" "batch:bool" "nocommandline" "requiredFields:string" "generateParamsStruct" "extends:string"

type GenOptsOption struct {
	f func(*genOptsOptionImpl)
	s string
}

func (o GenOptsOption) String() string { return o.s }

type GenOptsOptions interface {
	Batch() bool
	HasBatch() bool
	Extends() string
	HasExtends() bool
	Function() string
	HasFunction() bool
	GenerateParamsStruct() bool
	HasGenerateParamsStruct() bool
	Nocommandline() bool
	HasNocommandline() bool
	Outfile() string
	HasOutfile() bool
	Prefix() string
	HasPrefix() bool
	PrefixOptsType() bool
	HasPrefixOptsType() bool
	RequiredFields() string
	HasRequiredFields() bool
}

func GenOptsBatch(batch bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_batch = true
		opts.batch = batch
	}, fmt.Sprintf("gen.GenOptsBatch(bool %+v)}", batch)}
}
func GenOptsBatchFlag(batch *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if batch == nil {
			return
		}
		opts.has_batch = true
		opts.batch = *batch
	}, fmt.Sprintf("gen.GenOptsBatch(bool %+v)}", batch)}
}

func GenOptsExtends(extends string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_extends = true
		opts.extends = extends
	}, fmt.Sprintf("gen.GenOptsExtends(string %+v)}", extends)}
}
func GenOptsExtendsFlag(extends *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if extends == nil {
			return
		}
		opts.has_extends = true
		opts.extends = *extends
	}, fmt.Sprintf("gen.GenOptsExtends(string %+v)}", extends)}
}

func GenOptsFunction(function string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_function = true
		opts.function = function
	}, fmt.Sprintf("gen.GenOptsFunction(string %+v)}", function)}
}
func GenOptsFunctionFlag(function *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if function == nil {
			return
		}
		opts.has_function = true
		opts.function = *function
	}, fmt.Sprintf("gen.GenOptsFunction(string %+v)}", function)}
}

func GenOptsGenerateParamsStruct(generateParamsStruct bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_generateParamsStruct = true
		opts.generateParamsStruct = generateParamsStruct
	}, fmt.Sprintf("gen.GenOptsGenerateParamsStruct(bool %+v)}", generateParamsStruct)}
}
func GenOptsGenerateParamsStructFlag(generateParamsStruct *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if generateParamsStruct == nil {
			return
		}
		opts.has_generateParamsStruct = true
		opts.generateParamsStruct = *generateParamsStruct
	}, fmt.Sprintf("gen.GenOptsGenerateParamsStruct(bool %+v)}", generateParamsStruct)}
}

func GenOptsNocommandline(nocommandline bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_nocommandline = true
		opts.nocommandline = nocommandline
	}, fmt.Sprintf("gen.GenOptsNocommandline(bool %+v)}", nocommandline)}
}
func GenOptsNocommandlineFlag(nocommandline *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if nocommandline == nil {
			return
		}
		opts.has_nocommandline = true
		opts.nocommandline = *nocommandline
	}, fmt.Sprintf("gen.GenOptsNocommandline(bool %+v)}", nocommandline)}
}

func GenOptsOutfile(outfile string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_outfile = true
		opts.outfile = outfile
	}, fmt.Sprintf("gen.GenOptsOutfile(string %+v)}", outfile)}
}
func GenOptsOutfileFlag(outfile *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if outfile == nil {
			return
		}
		opts.has_outfile = true
		opts.outfile = *outfile
	}, fmt.Sprintf("gen.GenOptsOutfile(string %+v)}", outfile)}
}

func GenOptsPrefix(prefix string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_prefix = true
		opts.prefix = prefix
	}, fmt.Sprintf("gen.GenOptsPrefix(string %+v)}", prefix)}
}
func GenOptsPrefixFlag(prefix *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if prefix == nil {
			return
		}
		opts.has_prefix = true
		opts.prefix = *prefix
	}, fmt.Sprintf("gen.GenOptsPrefix(string %+v)}", prefix)}
}

func GenOptsPrefixOptsType(prefixOptsType bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_prefixOptsType = true
		opts.prefixOptsType = prefixOptsType
	}, fmt.Sprintf("gen.GenOptsPrefixOptsType(bool %+v)}", prefixOptsType)}
}
func GenOptsPrefixOptsTypeFlag(prefixOptsType *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if prefixOptsType == nil {
			return
		}
		opts.has_prefixOptsType = true
		opts.prefixOptsType = *prefixOptsType
	}, fmt.Sprintf("gen.GenOptsPrefixOptsType(bool %+v)}", prefixOptsType)}
}

func GenOptsRequiredFields(requiredFields string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_requiredFields = true
		opts.requiredFields = requiredFields
	}, fmt.Sprintf("gen.GenOptsRequiredFields(string %+v)}", requiredFields)}
}
func GenOptsRequiredFieldsFlag(requiredFields *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if requiredFields == nil {
			return
		}
		opts.has_requiredFields = true
		opts.requiredFields = *requiredFields
	}, fmt.Sprintf("gen.GenOptsRequiredFields(string %+v)}", requiredFields)}
}

type genOptsOptionImpl struct {
	prefixOptsType           bool
	has_prefixOptsType       bool
	prefix                   string
	has_prefix               bool
	function                 string
	has_function             bool
	outfile                  string
	has_outfile              bool
	batch                    bool
	has_batch                bool
	nocommandline            bool
	has_nocommandline        bool
	requiredFields           string
	has_requiredFields       bool
	generateParamsStruct     bool
	has_generateParamsStruct bool
	extends                  string
	has_extends              bool
}

func (g *genOptsOptionImpl) Batch() bool                   { return g.batch }
func (g *genOptsOptionImpl) HasBatch() bool                { return g.has_batch }
func (g *genOptsOptionImpl) Extends() string               { return g.extends }
func (g *genOptsOptionImpl) HasExtends() bool              { return g.has_extends }
func (g *genOptsOptionImpl) Function() string              { return g.function }
func (g *genOptsOptionImpl) HasFunction() bool             { return g.has_function }
func (g *genOptsOptionImpl) GenerateParamsStruct() bool    { return g.generateParamsStruct }
func (g *genOptsOptionImpl) HasGenerateParamsStruct() bool { return g.has_generateParamsStruct }
func (g *genOptsOptionImpl) Nocommandline() bool           { return g.nocommandline }
func (g *genOptsOptionImpl) HasNocommandline() bool        { return g.has_nocommandline }
func (g *genOptsOptionImpl) Outfile() string               { return g.outfile }
func (g *genOptsOptionImpl) HasOutfile() bool              { return g.has_outfile }
func (g *genOptsOptionImpl) Prefix() string                { return g.prefix }
func (g *genOptsOptionImpl) HasPrefix() bool               { return g.has_prefix }
func (g *genOptsOptionImpl) PrefixOptsType() bool          { return g.prefixOptsType }
func (g *genOptsOptionImpl) HasPrefixOptsType() bool       { return g.has_prefixOptsType }
func (g *genOptsOptionImpl) RequiredFields() string        { return g.requiredFields }
func (g *genOptsOptionImpl) HasRequiredFields() bool       { return g.has_requiredFields }

func makeGenOptsOptionImpl(opts ...GenOptsOption) *genOptsOptionImpl {
	res := &genOptsOptionImpl{}
	for _, opt := range opts {
		opt.f(res)
	}
	return res
}

func MakeGenOptsOptions(opts ...GenOptsOption) GenOptsOptions {
	return makeGenOptsOptionImpl(opts...)
}
