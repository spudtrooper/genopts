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

		% genopts --prefix Foo --outfile path/to/foooptions.go 'bar:bool' 'baz:int' 'boo:string'

2.  Add the new `float64` optiona by updating the commented command line in `path/to/foooptions.go` to be:

		// genopts --prefix Foo --outfile path/to/foooptions.go 'bar:bool' 'baz:int' 'boo:string' 'bam:float64'

3.  Update all relevant files in the current directory:

		% genopts --update

4.  Repeat

## Detailed usage

Get it:

```
go install github.com/spudtrooper/genopts
```

### Basics 

Run it:

```bash
~/go/bin/genopts --opts_type <type> <field-spec>+
```

Generates boilerplate for function objects named *type* with setters
for each *field-spec*, where each field is of the form `<name>*` 
for bool fields named *name* or `<name>:<type>` for fields named *name* 
and type *type*.

### Writing to files directly

To write to a file, pass `--outfile`, e.g.

```bash
~/go/bin/genopts --outfile=path/to/options.go --opts_type <type> <field-spec>+
```

### Batch

To update all the files under the current directory or the directory specified by `--update_dir`.

### Exclude directories in file search

Set `--exclude_dirs` to specify directories to exclude when searching for files with `--update`. This is useful if you have lots of non-go files under certain paths that will slow down the incremental runs of `genopts --update`.

### Config file

If you use `--update` to update files under a particular directory, you can also specify a *config* file that should contain a JSON-encoded version of a `Config` from `main.go`. By default we look for a file named `/.genopts` located in the directory specified by `--update_dir` or you can explicitly set this with `--config`, e.g.

```bash
~/go/bin/genopts --config path/to/config
```

The current options you can specify are `Excludes` and `GoImports` that will effectively set the flags `--exclude_dirs` and `--goimports`.

You can write the config with the `--write_config` flag. So, if you pass `--exclude_dirs` or `--goimports` everytime you run with `--update` and want to stop, you can write this file one and then stop passing these flags, e.g.

```bash
# Incremental updates with specific excluded directory and goimports path
genopts --update --excluded_dirs=foo,bar,baz --goimports=path/to/goimports
genopts --update --excluded_dirs=foo,bar,baz --goimports=path/to/goimports
genopts --update --excluded_dirs=foo,bar,baz --goimports=path/to/goimports
...

# Grow tired of passing --excluded_dirs and --goimports
# Decide to write the config
genopts --write_config \
        --update --excluded_dirs=foo,bar,baz --goimports=path/to/goimports

... writes the --excluded_dirs and --goimports flags to .genopts

# Subsequent calls --update will read the flags from .genopts
genopts --update
genopts --update
genopts --update
...
```
