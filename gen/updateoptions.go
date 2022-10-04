// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package gen

//go:generate genopts --prefix=Update --outfile=gen/updateoptions.go "threads:int"

type UpdateOption func(*updateOptionImpl)

type UpdateOptions interface {
	Threads() int
	HasThreads() bool
}

func UpdateThreads(threads int) UpdateOption {
	return func(opts *updateOptionImpl) {
		opts.has_threads = true
		opts.threads = threads
	}
}
func UpdateThreadsFlag(threads *int) UpdateOption {
	return func(opts *updateOptionImpl) {
		// if threads == nil {
		// 	return
		// }
		opts.has_threads = true
		opts.threads = *threads
	}
}

type updateOptionImpl struct {
	threads     int
	has_threads bool
}

func (u *updateOptionImpl) Threads() int     { return u.threads }
func (u *updateOptionImpl) HasThreads() bool { return u.has_threads }

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
