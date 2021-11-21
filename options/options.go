package options

// go run main.go prefixOptsType:bool prefix:string outfile:string

type Option func(*optionImpl)

type Options interface {
	PrefixOptsType() bool
	Prefix() string
	Outfile() string
}

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

func Outfile(outfile string) Option {
	return func(opts *optionImpl) {
		opts.outfile = outfile
	}
}

type optionImpl struct {
	prefixOptsType bool
	prefix         string
	outfile        string
}

func (o *optionImpl) PrefixOptsType() bool { return o.prefixOptsType }
func (o *optionImpl) Prefix() string       { return o.prefix }
func (o *optionImpl) Outfile() string      { return o.outfile }

func makeOptionImpl(opts ...Option) *optionImpl {
	res := &optionImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeOptions(opts ...Option) Options {
	return makeOptionImpl(opts...)
}
