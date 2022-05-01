package gen

//go:generate genopts --prefix=GenOpts --outfile=genoptsoptions.go "prefixOptsType:bool" "prefix:string" "function:string" "outfile:string" "batch:bool" "nocommandline"

type GenOptsOption func(*genOptsOptionImpl)

type GenOptsOptions interface {
	PrefixOptsType() bool
	Prefix() string
	Function() string
	Outfile() string
	Batch() bool
	Nocommandline() bool
}

func GenOptsPrefixOptsType(prefixOptsType bool) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.prefixOptsType = prefixOptsType
	}
}
func GenOptsPrefixOptsTypeFlag(prefixOptsType *bool) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.prefixOptsType = *prefixOptsType
	}
}

func GenOptsPrefix(prefix string) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.prefix = prefix
	}
}
func GenOptsPrefixFlag(prefix *string) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.prefix = *prefix
	}
}

func GenOptsFunction(function string) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.function = function
	}
}
func GenOptsFunctionFlag(function *string) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.function = *function
	}
}

func GenOptsOutfile(outfile string) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.outfile = outfile
	}
}
func GenOptsOutfileFlag(outfile *string) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.outfile = *outfile
	}
}

func GenOptsBatch(batch bool) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.batch = batch
	}
}
func GenOptsBatchFlag(batch *bool) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.batch = *batch
	}
}

func GenOptsNocommandline(nocommandline bool) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.nocommandline = nocommandline
	}
}
func GenOptsNocommandlineFlag(nocommandline *bool) GenOptsOption {
	return func(opts *genOptsOptionImpl) {
		opts.nocommandline = *nocommandline
	}
}

type genOptsOptionImpl struct {
	prefixOptsType bool
	prefix         string
	function       string
	outfile        string
	batch          bool
	nocommandline  bool
}

func (g *genOptsOptionImpl) PrefixOptsType() bool { return g.prefixOptsType }
func (g *genOptsOptionImpl) Prefix() string       { return g.prefix }
func (g *genOptsOptionImpl) Function() string     { return g.function }
func (g *genOptsOptionImpl) Outfile() string      { return g.outfile }
func (g *genOptsOptionImpl) Batch() bool          { return g.batch }
func (g *genOptsOptionImpl) Nocommandline() bool  { return g.nocommandline }

func makeGenOptsOptionImpl(opts ...GenOptsOption) *genOptsOptionImpl {
	res := &genOptsOptionImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeGenOptsOptions(opts ...GenOptsOption) GenOptsOptions {
	return makeGenOptsOptionImpl(opts...)
}
