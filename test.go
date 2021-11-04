package testdata

type SomeOption func(*someOptionImpl)

type SomeOptions interface {
	Foo() bool
	Bar() string
	Baz() float64
}

func SomeOptionFoo(foo bool) SomeOption {
	return func(opts *someOptionImpl) {
		opts.foo = foo
	}
}

func SomeOptionBar(bar string) SomeOption {
	return func(opts *someOptionImpl) {
		opts.bar = bar
	}
}

func SomeOptionBaz(baz float64) SomeOption {
	return func(opts *someOptionImpl) {
		opts.baz = baz
	}
}

type someOptionImpl struct {
	foo bool
	bar string
	baz float64
}

func (s *someOptionImpl) Foo() bool    { return s.foo }
func (s *someOptionImpl) Bar() string  { return s.bar }
func (s *someOptionImpl) Baz() float64 { return s.baz }

func makeSomeOptionImpl(opts ...SomeOption) someOptionImpl {
	var res someOptionImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
