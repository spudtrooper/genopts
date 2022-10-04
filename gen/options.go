// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package gen

//go:generate genopts --opt_type=SomeOption --outfile=gen/options.go

type SomeOption func(*someOptionImpl)

type SomeOptions interface {
}

type someOptionImpl struct {
}

func makeSomeOptionImpl(opts ...SomeOption) *someOptionImpl {
	res := &someOptionImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeSomeOptions(opts ...SomeOption) SomeOptions {
	return makeSomeOptionImpl(opts...)
}
