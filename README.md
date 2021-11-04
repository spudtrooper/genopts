# genopts

A command line tool to generate "options" functions in go.

## Usage

Get it:

```
go install github.com/spudtrooper/genopts
```

Run it:

```
~/go/bin/genopts --opts_type <type> <field-spec>+
```

Generates boilerplate for function objects named *type* with setters
for each *field-spec*, where each field is of the form `<name>*` 
for bool fields named *name* or `<name>:<type>` for fields named *name* 
and type *type*.

## Example


Generate your options code with:

```
genopts --opt_type SomeOpt foo bar:string baz:float64
```

...then you could use this generated code with something like:

```go
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
```

The generated code would look like this:

```go
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

func makeSomeOptImpl(opts ...SomeOpt) someOptImpl {
	var res someOptImpl
	for _, opt := range opts {
		opt(&res)
	}
	return res
}

func MakeSomeOpts(opts ...SomeOpt) SomeOpts {
	return makeSomeOptImpl(opts...)
}
```

