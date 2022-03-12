package driver

import "flag"

var (
	boolFlag   = flag.Bool("bool", false, "some bool flag")
	stringFlag = flag.String("string", "", "some string flag")
	floatFlag  = flag.Float64("float", 0, "some float flag")
)

func Func() {
	TakesOpts(Foo(true), Bar("bar"), Baz(1.0))
	TakesOpts(FooFlag(boolFlag), Bar(stringFlag), Baz(floatFlag))
}

func TakesOpts(opts ...SomeOpts) {
	// nothing
}
