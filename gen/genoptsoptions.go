package gen

//go:generate genopts --prefix=GenOpts --outfile=gen/genoptsoptions.go "prefixOptsType:bool" "prefix:string" "outfile:string" "batch:bool"

type GenOptsOption func(*genOptsOptionImpl)

type GenOptsOptions interface {
	PrefixOptsType() bool
	Prefix() string
	Outfile() string
	Batch() bool
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

type genOptsOptionImpl struct {
	prefixOptsType bool
	prefix         string
	outfile        string
	batch          bool
}

func (g *genOptsOptionImpl) PrefixOptsType() bool { return g.prefixOptsType }
func (g *genOptsOptionImpl) Prefix() string       { return g.prefix }
func (g *genOptsOptionImpl) Outfile() string      { return g.outfile }
func (g *genOptsOptionImpl) Batch() bool          { return g.batch }

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
