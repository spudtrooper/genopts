package genopts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

func TestGenOpts(t *testing.T) {
	var tests = []struct {
		name     string
		optsType string
		implType string
		fields   []string
		opts     []Option
		want     string
	}{
		{
			name:     "empty",
			optsType: "SomeOptions",
			implType: "",
			fields:   []string{},
			want: `
type SomeOptions func(*someOptionsImpl)

type someOptionsImpl struct {

}

func makeSomeOptionsImpl(opts ...SomeOptions) someOptionsImpl {
	var res someOptionsImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`,
		},
		{
			name:     "impl",
			optsType: "SomeOptions",
			implType: "explicitImpl",
			fields:   []string{},
			want: `
type SomeOptions func(*explicitImpl)

type explicitImpl struct {

}

func makeExplicitImpl(opts ...SomeOptions) explicitImpl {
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
			optsType: "SomeOptions",
			implType: "",
			fields:   []string{"foo", "bar:string", "baz:float64"},
			want: `
type SomeOptions func(*someOptionsImpl)

func Foo(foo bool) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.foo = foo
	}
}

func Bar(bar string) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.bar = bar
	}
}

func Baz(baz float64) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.baz = baz
	}
}

type someOptionsImpl struct {
	foo bool
	bar string
	baz float64

}

func makeSomeOptionsImpl(opts ...SomeOptions) someOptionsImpl {
	var res someOptionsImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`,
		}, {
			name:     "prefix",
			optsType: "SomeOptions",
			implType: "",
			fields:   []string{"foo", "bar:string", "baz:float64"},
			opts: []Option{
				Prefix("Prefix"),
			},
			want: `
type SomeOptions func(*someOptionsImpl)

func PrefixFoo(foo bool) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.foo = foo
	}
}

func PrefixBar(bar string) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.bar = bar
	}
}

func PrefixBaz(baz float64) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.baz = baz
	}
}

type someOptionsImpl struct {
	foo bool
	bar string
	baz float64

}

func makeSomeOptionsImpl(opts ...SomeOptions) someOptionsImpl {
	var res someOptionsImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
`,
		}, {
			name:     "prefixOptsType",
			optsType: "SomeOptions",
			implType: "",
			fields:   []string{"foo", "bar:string", "baz:float64"},
			opts: []Option{
				PrefixOptsType(true),
			},
			want: `
type SomeOptions func(*someOptionsImpl)

func SomeOptionsFoo(foo bool) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.foo = foo
	}
}

func SomeOptionsBar(bar string) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.bar = bar
	}
}

func SomeOptionsBaz(baz float64) SomeOptions {
	return func(opts *someOptionsImpl) {
		opts.baz = baz
	}
}

type someOptionsImpl struct {
	foo bool
	bar string
	baz float64

}

func makeSomeOptionsImpl(opts ...SomeOptions) someOptionsImpl {
	var res someOptionsImpl
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
			got, err := GenOpts(test.optsType, test.implType, test.fields, test.opts...)
			if err != nil {
				t.Fatalf("GenOpts(%q,%q,%v): %v", test.optsType, test.implType, test.fields, err)
			}
			if want, got := strings.TrimSpace(test.want), strings.TrimSpace(got); want != got {
				fmt.Println(got)
				t.Errorf("GenOpts(%q,%q,%v) want != got:\n\n------\n%s\n-------", test.optsType, test.implType, test.fields, diff.LineDiff(want, got))
			}
		})
	}
}
