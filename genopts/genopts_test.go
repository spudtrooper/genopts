package genopts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/spudtrooper/genopts/options"
)

func TestGenOpts(t *testing.T) {
	var tests = []struct {
		name     string
		optType  string
		implType string
		fields   []string
		opts     []options.Option
		want     string
	}{
		{
			name:     "empty",
			optType:  "SomeOption",
			implType: "",
			fields:   []string{},
			want: `
type SomeOption func(*someOptionImpl)

type SomeOptions interface {

}

type someOptionImpl struct {

}


func makeSomeOptionImpl(opts ...SomeOption) someOptionImpl {
	var res someOptionImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`,
		},
		{
			name:     "impl",
			optType:  "SomeOption",
			implType: "explicitImpl",
			fields:   []string{},
			want: `
type SomeOption func(*explicitImpl)

type SomeOptions interface {

}

type explicitImpl struct {

}


func makeExplicitImpl(opts ...SomeOption) explicitImpl {
	var res explicitImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`,
		},
		{
			name:     "fields",
			optType:  "SomeOption",
			implType: "",
			fields:   []string{"foo", "bar:string", "baz:float64"},
			want: `
type SomeOption func(*someOptionImpl)

type SomeOptions interface {
	Foo() bool
	Bar() string
	Baz() float64

}

func Foo(foo bool) SomeOption {
	return func(opts *someOptionImpl) {
		opts.foo = foo
	}
}

func Bar(bar string) SomeOption {
	return func(opts *someOptionImpl) {
		opts.bar = bar
	}
}

func Baz(baz float64) SomeOption {
	return func(opts *someOptionImpl) {
		opts.baz = baz
	}
}

type someOptionImpl struct {
	foo bool
	bar string
	baz float64

}

func (s *someOptionImpl) Foo() bool { return s.foo }
func (s *someOptionImpl) Bar() string { return s.bar }
func (s *someOptionImpl) Baz() float64 { return s.baz }

func makeSomeOptionImpl(opts ...SomeOption) someOptionImpl {
	var res someOptionImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`,
		}, {
			name:     "prefix",
			optType:  "SomeOption",
			implType: "",
			fields:   []string{"foo", "bar:string", "baz:float64"},
			opts: []options.Option{
				options.Prefix("Prefix"),
			},
			want: `
type SomeOption func(*someOptionImpl)

type SomeOptions interface {
	Foo() bool
	Bar() string
	Baz() float64

}

func PrefixFoo(foo bool) SomeOption {
	return func(opts *someOptionImpl) {
		opts.foo = foo
	}
}

func PrefixBar(bar string) SomeOption {
	return func(opts *someOptionImpl) {
		opts.bar = bar
	}
}

func PrefixBaz(baz float64) SomeOption {
	return func(opts *someOptionImpl) {
		opts.baz = baz
	}
}

type someOptionImpl struct {
	foo bool
	bar string
	baz float64

}

func (s *someOptionImpl) Foo() bool { return s.foo }
func (s *someOptionImpl) Bar() string { return s.bar }
func (s *someOptionImpl) Baz() float64 { return s.baz }

func makeSomeOptionImpl(opts ...SomeOption) someOptionImpl {
	var res someOptionImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`,
		}, {
			name:     "prefixOptsType",
			optType:  "SomeOption",
			implType: "",
			fields:   []string{"foo", "bar:string", "baz:float64"},
			opts: []options.Option{
				options.PrefixOptsType(true),
			},
			want: `
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

func (s *someOptionImpl) Foo() bool { return s.foo }
func (s *someOptionImpl) Bar() string { return s.bar }
func (s *someOptionImpl) Baz() float64 { return s.baz }

func makeSomeOptionImpl(opts ...SomeOption) someOptionImpl {
	var res someOptionImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GenOpts(test.optType, test.implType, test.fields, test.opts...)
			if err != nil {
				t.Fatalf("GenOpts(%q,%q,%v): %v", test.optType, test.implType, test.fields, err)
			}
			if false {
				fmt.Printf(" >>>>>>>>>>>>>>> GOT (%s) >>>>>>>>>>>>>>>>>>>>>>\n", test.name)
				fmt.Println(got)
				fmt.Printf(" <<<<<<<<<<<<<<< GOT (%s) <<<<<<<<<<<<<<<<<<<<<<\n", test.name)
			}
			if want, got := strings.TrimSpace(test.want), strings.TrimSpace(got); want != got {
				fmt.Println(got)
				t.Errorf("GenOpts(%q,%q,%v) want != got:\n\n------\n%s\n-------", test.optType, test.implType, test.fields, diff.LineDiff(want, got))
			}
		})
	}
}
