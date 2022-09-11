# genopts

A go generator for *optional* parameters to functions.

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

you can declaratively generate code to handle these optional parameters with `//go:generate`:

```bash
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

## Installation

Get it:

```
go install github.com/spudtrooper/genopts@latest
```

and make sure the binary is in your path, e.g. add the following to your .bashrc or .zprofile:

```bash
export PATH=$PATH:~/go/bin
```

### Usage 

Either run as a go generator with something like

```bash
//go:generate genopts --opts_type <type> <field-spec>+ --outfile=<this file>
```

or run it directly:

```bash
~/go/bin/genopts --opts_type <type> <field-spec>+
```

This generates boilerplate for function objects named *type* with setters
for each *field-spec*, where each field is of the form `<name>*` 
for bool fields named *name* or `<name>:<type>` for fields named *name* 
and type *type*.

### Writing to files directly

To write to a file, pass `--outfile`, e.g.

```bash
~/go/bin/genopts --outfile=path/to/options.go --opts_type <type> <field-spec>+
```

### Batch

To update all the files under the current directory run either

```
go generate ./...
```

or 

```
genopts
```

### Exclude directories in file search

Set `--exclude_dirs` with a comma-delimited list of directories to specify which directories to exclude when searching for files with `--update`. This is useful if you have lots of non-go files under certain paths that will slow down the incremental runs of `genopts --update`.

### Config file

If you use `--update` to update files under a particular directory, you can also specify a *config* file that should contain a JSON-encoded version of a `Config` from `main.go`. By default we look for a file named `/.genopts` located in the directory specified by `--update_dir` or you can explicitly set this with `--config`.

The current options you can specify are `Excludes` and `GoImports` that will effectively set the flags `--exclude_dirs` and `--goimports`.

#### Config file idiom

You can write the config with the `--write_config` flag. So, if you pass `--exclude_dirs` or `--goimports` everytime you run with `--update` and want to stop, you can write this file one and then stop passing these flags, e.g.

```bash
# Incremental updates with specific excluded directory and goimports path
genopts --update --excluded_dirs=foo,bar,baz --goimports=path/to/goimports
genopts --update --excluded_dirs=foo,bar,baz --goimports=path/to/goimports
genopts --update --excluded_dirs=foo,bar,baz --goimports=path/to/goimports
...

# Grow tired of passing --excluded_dirs and --goimports, you conclude that the only option is 
# to end it all or write the config. You decide to write to the config.
genopts --update --excluded_dirs=foo,bar,baz --goimports=path/to/goimports --write_config

... writes the --excluded_dirs and --goimports flags to .genopts

# Subsequent calls to --update will read the flags from .genopts
genopts --update
genopts --update
genopts --update
...
```
