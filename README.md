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

```
genopts --opts_type SomeOpts foo bar:string baz:float64
```

generates something like:

```go
// SomeOpts are options to TODO
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
```

