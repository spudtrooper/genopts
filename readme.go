package readme

import (
	"flag"
	"fmt"
)

type SomeOpt func(*someOptImpl)

type SomeOpts interface {
	Foo() bool
	Bar() string
	Baz() float64
}

func Foo(foo bool) SomeOpt {
	return func(opts *someOptImpl) {
		opts.foo = foo
	}
}

func Bar(bar string) SomeOpt {
	return func(opts *someOptImpl) {
		opts.bar = bar
	}
}

func Baz(baz float64) SomeOpt {
	return func(opts *someOptImpl) {
		opts.baz = baz
	}
}

type someOptImpl struct {
	foo bool
	bar string
	baz float64
}

func (s *someOptImpl) Foo() bool    { return s.foo }
func (s *someOptImpl) Bar() string  { return s.bar }
func (s *someOptImpl) Baz() float64 { return s.baz }

func makeSomeOptImpl(opts ...SomeOpt) *someOptImpl {
	res := &someOptImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeSomeOpts(opts ...SomeOpt) SomeOpts {
	return makeSomeOptImpl(opts...)
}

var (
	foo = flag.Bool("foo", false, "some bool flag")
	bar = flag.String("bar", "", "some string flag")
	baz = flag.Float64("baz", 0, "some float64 flag")
)

func consumesOptions(inputOpts ...SomeOpt) {
	opts := MakeSomeOpts(inputOpts...)
	if opts.Foo() {
		fmt.Println(opts.Bar())
	}
}

func producesOptions() {
	consumesOptions(Foo(*foo), Bar(*bar), Baz(*baz))
}

func main() {

}
