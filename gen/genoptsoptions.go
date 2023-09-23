// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package gen

import "fmt"

//go:generate genopts --prefix=GenOpts --outfile=genoptsoptions.go "prefixOptsType:bool" "prefix:string" "function:string" "outfile:string" "batch:bool" "nocommandline" "requiredFields:string" "generateParamsStruct" "extends:string" "musts"

type GenOptsOption struct {
	f func(*genOptsOptionImpl)
	s string
}

func (o GenOptsOption) String() string { return o.s }

type GenOptsOptions interface {
	Batch() bool
	HasBatch() bool
	MustBatch() bool
	Extends() string
	HasExtends() bool
	MustExtends() string
	Function() string
	HasFunction() bool
	MustFunction() string
	GenerateParamsStruct() bool
	HasGenerateParamsStruct() bool
	MustGenerateParamsStruct() bool
	Musts() bool
	HasMusts() bool
	MustMusts() bool
	Nocommandline() bool
	HasNocommandline() bool
	MustNocommandline() bool
	Outfile() string
	HasOutfile() bool
	MustOutfile() string
	Prefix() string
	HasPrefix() bool
	MustPrefix() string
	PrefixOptsType() bool
	HasPrefixOptsType() bool
	MustPrefixOptsType() bool
	RequiredFields() string
	HasRequiredFields() bool
	MustRequiredFields() string
}

func GenOptsBatch(batch bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_batch = true
		opts.batch = batch
	}, fmt.Sprintf("gen.GenOptsBatch(bool %+v)", batch)}
}
func GenOptsBatchFlag(batch *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if batch == nil {
			return
		}
		opts.has_batch = true
		opts.batch = *batch
	}, fmt.Sprintf("gen.GenOptsBatch(bool %+v)", batch)}
}

func GenOptsExtends(extends string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_extends = true
		opts.extends = extends
	}, fmt.Sprintf("gen.GenOptsExtends(string %+v)", extends)}
}
func GenOptsExtendsFlag(extends *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if extends == nil {
			return
		}
		opts.has_extends = true
		opts.extends = *extends
	}, fmt.Sprintf("gen.GenOptsExtends(string %+v)", extends)}
}

func GenOptsFunction(function string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_function = true
		opts.function = function
	}, fmt.Sprintf("gen.GenOptsFunction(string %+v)", function)}
}
func GenOptsFunctionFlag(function *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if function == nil {
			return
		}
		opts.has_function = true
		opts.function = *function
	}, fmt.Sprintf("gen.GenOptsFunction(string %+v)", function)}
}

func GenOptsGenerateParamsStruct(generateParamsStruct bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_generateParamsStruct = true
		opts.generateParamsStruct = generateParamsStruct
	}, fmt.Sprintf("gen.GenOptsGenerateParamsStruct(bool %+v)", generateParamsStruct)}
}
func GenOptsGenerateParamsStructFlag(generateParamsStruct *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if generateParamsStruct == nil {
			return
		}
		opts.has_generateParamsStruct = true
		opts.generateParamsStruct = *generateParamsStruct
	}, fmt.Sprintf("gen.GenOptsGenerateParamsStruct(bool %+v)", generateParamsStruct)}
}

func GenOptsMusts(musts bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_musts = true
		opts.musts = musts
	}, fmt.Sprintf("gen.GenOptsMusts(bool %+v)", musts)}
}
func GenOptsMustsFlag(musts *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if musts == nil {
			return
		}
		opts.has_musts = true
		opts.musts = *musts
	}, fmt.Sprintf("gen.GenOptsMusts(bool %+v)", musts)}
}

func GenOptsNocommandline(nocommandline bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_nocommandline = true
		opts.nocommandline = nocommandline
	}, fmt.Sprintf("gen.GenOptsNocommandline(bool %+v)", nocommandline)}
}
func GenOptsNocommandlineFlag(nocommandline *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if nocommandline == nil {
			return
		}
		opts.has_nocommandline = true
		opts.nocommandline = *nocommandline
	}, fmt.Sprintf("gen.GenOptsNocommandline(bool %+v)", nocommandline)}
}

func GenOptsOutfile(outfile string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_outfile = true
		opts.outfile = outfile
	}, fmt.Sprintf("gen.GenOptsOutfile(string %+v)", outfile)}
}
func GenOptsOutfileFlag(outfile *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if outfile == nil {
			return
		}
		opts.has_outfile = true
		opts.outfile = *outfile
	}, fmt.Sprintf("gen.GenOptsOutfile(string %+v)", outfile)}
}

func GenOptsPrefix(prefix string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_prefix = true
		opts.prefix = prefix
	}, fmt.Sprintf("gen.GenOptsPrefix(string %+v)", prefix)}
}
func GenOptsPrefixFlag(prefix *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if prefix == nil {
			return
		}
		opts.has_prefix = true
		opts.prefix = *prefix
	}, fmt.Sprintf("gen.GenOptsPrefix(string %+v)", prefix)}
}

func GenOptsPrefixOptsType(prefixOptsType bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_prefixOptsType = true
		opts.prefixOptsType = prefixOptsType
	}, fmt.Sprintf("gen.GenOptsPrefixOptsType(bool %+v)", prefixOptsType)}
}
func GenOptsPrefixOptsTypeFlag(prefixOptsType *bool) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if prefixOptsType == nil {
			return
		}
		opts.has_prefixOptsType = true
		opts.prefixOptsType = *prefixOptsType
	}, fmt.Sprintf("gen.GenOptsPrefixOptsType(bool %+v)", prefixOptsType)}
}

func GenOptsRequiredFields(requiredFields string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		opts.has_requiredFields = true
		opts.requiredFields = requiredFields
	}, fmt.Sprintf("gen.GenOptsRequiredFields(string %+v)", requiredFields)}
}
func GenOptsRequiredFieldsFlag(requiredFields *string) GenOptsOption {
	return GenOptsOption{func(opts *genOptsOptionImpl) {
		if requiredFields == nil {
			return
		}
		opts.has_requiredFields = true
		opts.requiredFields = *requiredFields
	}, fmt.Sprintf("gen.GenOptsRequiredFields(string %+v)", requiredFields)}
}

type genOptsOptionImpl struct {
	batch                    bool
	has_batch                bool
	extends                  string
	has_extends              bool
	function                 string
	has_function             bool
	generateParamsStruct     bool
	has_generateParamsStruct bool
	musts                    bool
	has_musts                bool
	nocommandline            bool
	has_nocommandline        bool
	outfile                  string
	has_outfile              bool
	prefix                   string
	has_prefix               bool
	prefixOptsType           bool
	has_prefixOptsType       bool
	requiredFields           string
	has_requiredFields       bool
}

func (g *genOptsOptionImpl) Batch() bool    { return g.batch }
func (g *genOptsOptionImpl) HasBatch() bool { return g.has_batch }
func (g *genOptsOptionImpl) MustBatch() bool {
	if !g.has_batch {
		panic("batch is required")
	}
	return g.batch
}
func (g *genOptsOptionImpl) Extends() string  { return g.extends }
func (g *genOptsOptionImpl) HasExtends() bool { return g.has_extends }
func (g *genOptsOptionImpl) MustExtends() string {
	if !g.has_extends {
		panic("extends is required")
	}
	return g.extends
}
func (g *genOptsOptionImpl) Function() string  { return g.function }
func (g *genOptsOptionImpl) HasFunction() bool { return g.has_function }
func (g *genOptsOptionImpl) MustFunction() string {
	if !g.has_function {
		panic("function is required")
	}
	return g.function
}
func (g *genOptsOptionImpl) GenerateParamsStruct() bool    { return g.generateParamsStruct }
func (g *genOptsOptionImpl) HasGenerateParamsStruct() bool { return g.has_generateParamsStruct }
func (g *genOptsOptionImpl) MustGenerateParamsStruct() bool {
	if !g.has_generateParamsStruct {
		panic("generateParamsStruct is required")
	}
	return g.generateParamsStruct
}
func (g *genOptsOptionImpl) Musts() bool    { return g.musts }
func (g *genOptsOptionImpl) HasMusts() bool { return g.has_musts }
func (g *genOptsOptionImpl) MustMusts() bool {
	if !g.has_musts {
		panic("musts is required")
	}
	return g.musts
}
func (g *genOptsOptionImpl) Nocommandline() bool    { return g.nocommandline }
func (g *genOptsOptionImpl) HasNocommandline() bool { return g.has_nocommandline }
func (g *genOptsOptionImpl) MustNocommandline() bool {
	if !g.has_nocommandline {
		panic("nocommandline is required")
	}
	return g.nocommandline
}
func (g *genOptsOptionImpl) Outfile() string  { return g.outfile }
func (g *genOptsOptionImpl) HasOutfile() bool { return g.has_outfile }
func (g *genOptsOptionImpl) MustOutfile() string {
	if !g.has_outfile {
		panic("outfile is required")
	}
	return g.outfile
}
func (g *genOptsOptionImpl) Prefix() string  { return g.prefix }
func (g *genOptsOptionImpl) HasPrefix() bool { return g.has_prefix }
func (g *genOptsOptionImpl) MustPrefix() string {
	if !g.has_prefix {
		panic("prefix is required")
	}
	return g.prefix
}
func (g *genOptsOptionImpl) PrefixOptsType() bool    { return g.prefixOptsType }
func (g *genOptsOptionImpl) HasPrefixOptsType() bool { return g.has_prefixOptsType }
func (g *genOptsOptionImpl) MustPrefixOptsType() bool {
	if !g.has_prefixOptsType {
		panic("prefixOptsType is required")
	}
	return g.prefixOptsType
}
func (g *genOptsOptionImpl) RequiredFields() string  { return g.requiredFields }
func (g *genOptsOptionImpl) HasRequiredFields() bool { return g.has_requiredFields }
func (g *genOptsOptionImpl) MustRequiredFields() string {
	if !g.has_requiredFields {
		panic("requiredFields is required")
	}
	return g.requiredFields
}

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
