package some



type SomeOpts func(*someOptsImpl)

func Foo(foo bool) SomeOpts {
	return func(opts *someOptsImpl) {
		opts.foo = foo
	}
}

func Bar(bar string) SomeOpts {
	return func(opts *someOptsImpl) {
		opts.bar = bar
	}
}

func Baz(baz float64) SomeOpts {
	return func(opts *someOptsImpl) {
		opts.baz = baz
	}
}

type someOptsImpl struct {
	foo bool
	bar string
	baz float64

}
