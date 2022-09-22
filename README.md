# genopts

A go generator for *optional* parameters to functions.

## tl;dr

Instead of having functions with explicit *optional* parameters like:

```go
func Foo(requiredString string, optBar bool, optBaz int, optBoo string) {
  if optBar == false {
    ...
  } else {
    ...
  }
  if optBaz == 0 {
    ...
  } else {
    ...
  }
  if optBoo == "" {
    ...
  } else {
    ...
  }
}

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

you can declaratively generate code to handle these optional parameters with `//go:generate`:

```go
//go:generate genopts --function Foo 'bar:bool' 'baz:int' 'boo:string'
func Foo(...) {}
```

so you could have cleaner and easier-to-maintain code like:

```go
import "github.com/spudtrooper/goutil/or"

//go:generate genopts --function Foo 'bar:bool' 'baz:int' 'boo:string'
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

See some [real examples](https://github.com/search?q=%22go%3Agenerate+genopts%22&type=code).

## Installation

Get it:

```
go install github.com/spudtrooper/genopts@latest
```

and make sure the binary is in your path, e.g. add the following to your .bashrc or .zprofile:

```bash
export PATH=$PATH:~/go/bin
```

### Details 

See [details](https://github.com/spudtrooper/genopts/blob/main/details.md).