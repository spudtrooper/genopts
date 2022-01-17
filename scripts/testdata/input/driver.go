package driver

func Func() {
	TakesOpts(Foo(true), Bar("bar"), Baz(1.0))
}

func TakesOpts(opts ...SomeOpts) {
	// nothing
}
