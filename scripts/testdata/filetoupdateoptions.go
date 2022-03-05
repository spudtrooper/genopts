package testdata

//go:generate genopts --opt_type=FileToUpdateOption --prefix=FileToUpdate --outfile=some/filetoupdate.go "foo" "bar:int" "baz:string"

type FileToUpdateOption func(*fileToUpdateOptionImpl)

type FileToUpdateOptions interface {
	Foo() bool
	Bar() int
	Baz() string
}

func FileToUpdateFoo(foo bool) FileToUpdateOption {
	return func(opts *fileToUpdateOptionImpl) {
		opts.foo = foo
	}
}

func FileToUpdateBar(bar int) FileToUpdateOption {
	return func(opts *fileToUpdateOptionImpl) {
		opts.bar = bar
	}
}

func FileToUpdateBaz(baz string) FileToUpdateOption {
	return func(opts *fileToUpdateOptionImpl) {
		opts.baz = baz
	}
}

type fileToUpdateOptionImpl struct {
	foo bool
	bar int
	baz string
}

func (f *fileToUpdateOptionImpl) Foo() bool   { return f.foo }
func (f *fileToUpdateOptionImpl) Bar() int    { return f.bar }
func (f *fileToUpdateOptionImpl) Baz() string { return f.baz }

func makeFileToUpdateOptionImpl(opts ...FileToUpdateOption) *fileToUpdateOptionImpl {
	res := &fileToUpdateOptionImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeFileToUpdateOptions(opts ...FileToUpdateOption) FileToUpdateOptions {
	return makeFileToUpdateOptionImpl(opts...)
}
