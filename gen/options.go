// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package gen

//go:generate genopts --outfile=options.go

type Option func(*optionImpl)

type Options interface {
}

type optionImpl struct {
}

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