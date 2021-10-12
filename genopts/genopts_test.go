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
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GenOpts(test.optsType, test.implType, test.fields)
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
