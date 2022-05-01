// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package gen

//go:generate genopts --prefix=Update --outfile=updateoptions.go "threads:int"

type UpdateOption func(*updateOptionImpl)

type UpdateOptions interface {
	Threads() int
}

func UpdateThreads(threads int) UpdateOption {
	return func(opts *updateOptionImpl) {
		opts.threads = threads
	}
}
func UpdateThreadsFlag(threads *int) UpdateOption {
	return func(opts *updateOptionImpl) {
		opts.threads = *threads
	}
}

type updateOptionImpl struct {
	threads int
}

func (u *updateOptionImpl) Threads() int { return u.threads }

func makeUpdateOptionImpl(opts ...UpdateOption) *updateOptionImpl {
	res := &updateOptionImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeUpdateOptions(opts ...UpdateOption) UpdateOptions {
	return makeUpdateOptionImpl(opts...)
}
