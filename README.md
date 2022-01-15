# genopts

A command line tool to generate *optional* parameters for Go functions.

## tl;dr

Instead of having functions with explicit *optional* parameters like:

```go
func Foo(requiredString string, optBar bool, optBaz int, optBoo string) { ... }

func Usage() {
	...
	Foo("some required 1", true, 0, "")
	...
	Foo("some required 2", false, 1, "")
	...
	Foo("some required 3", false, 0, "optional string")
	...
}
```

you can run:

```bash
genopts --prefix Foo 'bar:bool' 'baz:int' 'boo:string'
```

which will generate the interface `FooOption`, constructor `MakeFooOptions`, and wrappers `FooBar(bool)`, `FooBaz(int)` & `FooBoo(string)`, so you could have cleaner and easier-to-maintain code like:

```go
import "github.com/spudtrooper/goutil/or"

func Foo(requiredString string bool, fOpts...FooOption) { 
	opts := MakeFooOptions(fOpts...)

	bar := or.Bool(opts.Bar(), false)
	baz := or.Int(opts.Baz(), 0)
	boo := or.String(opts.Boo(), "")

	...
 }

func Usage() {
	...
	Foo("some required 1", FooBar(true))
	...
	Foo("some required 2", FooBaz(1))
	...
	Foo("some required 3", FooBoo("optional string"))
	...
}
```

## Motivation

Say you have a function that takes one required arguments and three optional arguments, like:

```go
func Foo(requiredString string bool, optInt int, optString string) { ... }
```

The pain of treating the "optional" arguments as optional is:

1.	All callers have to pass these, and
1.  All callers have to by sync w.r.t. the **correct defaults**

(2) can be alleviated by always passing canonical defaults (e.g. `false` bools, `0` numbers, `""` strings), but there's still room for error. 

Assuming we only have the first problem, if we have 3 calls to `Foo`, each passing the required string and one of the optional arguments then we would have the following:

```go
func Foo(requiredString string, optBar bool, optBaz int, optBoo string) { ... }

func Usage() {
	...
	Foo("some required 1", true, 0, "")
	...
	Foo("some required 2", false, 1, "")
	...
	Foo("some required 3", false, 0, "optional string")
	...
}
```

Some corrolaries to the pain above are:

1. If you add an parameter to `Foo`, you have to update each caller, e.g.

```go
func Foo(requiredString string, optBar bool, optBaz int, optBoo string, newOptBam float64) { ... }

func Usage() {
	...
	Foo("some required 1", true, 0, "", 0)
	...
	Foo("some required 2", false, 1, "", 0)
	...
	Foo("some required 3", false, 0, "optional string", 0)
	...
}
```

2. If you have multiple parameters of the same type it can get confusing, e.g.

```go
func Foo2(requiredString string bool, optInt1 int, optInt2 int, optInt3 int) { ... }

func Usage() {
	...
	Foo2("some required 1", 1, 0, 0)
	...
	Foo2("some required 2", 0, 2, 0)
	...
	Foo2("some required 3", 0, 0, 3)
	...
}
```

## Solution

The solution provided by this package is easily-generated code for named optional parameters to go functions. So, the above example would look like this instead:

```go
import "github.com/spudtrooper/goutil/or"

func Foo(requiredString string bool, fOpts...FooOption) { 
	opts := MakeFooOptions(fOpts...)

	bar := or.Bool(opts.Bar(), false)
	baz := or.Int(opts.Bar(), 0)
	boo := or.String(opts.Bar(), "")

	...
 }

func Usage() {
	...
	Foo("some required 1", FooBar(true))
	...
	Foo("some required 2", FooBaz(1))
	...
	Foo("some required 3", FooBoo("optional string"))
	...
}
```

To generate the code for this you would run:

```bash
genopts --prefix Foo 'bar:bool' 'baz:int' 'boo:string'
```

The benefits are:

1.  Adding another parameter is easy, you just rerun the command above with one more argument.
1.  You don't have to update the existing callers
1.  The parameters are named, so you don't run into the many-parameters-of-the-same-type problem.
1.  Defaults are controlled by the functions; callers don't have to be in sync with the defaults
1.  Caller code is cleaner; i.e. is clear what non-defaults are being passed because these are the **only** arguments passed.

## Idiom

1.  Have one file per function
2.  Use the `--outfile` flag to write directly to the file. This file will contain the command line that generated the file as comments.
3.  When you want to modify existing options, modify the commented command line in the file and update with:

		% genopts --update


e.g. from above:

1. Generate the initial options

		% genopts --prefix Foo --outfile path/to/foooptions.go \
			'bar:bool' 'baz:int' 'boo:string'

2.  To add the new `float64` optiona, update the commented command line in `path/to/foooptions.go` to be:

		// genopts --prefix Foo --outfile path/to/foooptions.go \
			'bar:bool' 'baz:int' 'boo:string' 'bam:float64'

3.  Update all relevant files in the current directory:

		% genopts --update

## Usage

Get it:

```
go install github.com/spudtrooper/genopts
```

Run it:

```bash
~/go/bin/genopts --opts_type <type> <field-spec>+
```

Generates boilerplate for function objects named *type* with setters
for each *field-spec*, where each field is of the form `<name>*` 
for bool fields named *name* or `<name>:<type>` for fields named *name* 
and type *type*.

Run in batch:

```
~/go/bin/genopts --update
```

To update all the files under the current directory or the directory specified by `--update_dir`.

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

