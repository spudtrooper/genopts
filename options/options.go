package options

// Generate with:
// go run main.go prefixOptsType:bool  prefix:string
type Option func(*optionImpl)

type Options interface {
	PrefixOptsType() bool
	Prefix() string
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

type optionImpl struct {
	prefixOptsType bool
	prefix         string
}

func (o *optionImpl) PrefixOptsType() bool { return o.prefixOptsType }
func (o *optionImpl) Prefix() string       { return o.prefix }

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
