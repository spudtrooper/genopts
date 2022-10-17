// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package gen

//go:generate genopts --prefix=Update --outfile=gen/updateoptions.go "threads:int"

type UpdateOption struct {
	f func(*updateOptionImpl)
	s string
}

func (o UpdateOption) String() string { return o.s }

type UpdateOptions interface {
	Threads() int
	HasThreads() bool
}

func UpdateThreads(threads int) UpdateOption {
	return UpdateOption{func(opts *updateOptionImpl) {
		opts.has_threads = true
		opts.threads = threads
	}, "gen.UpdateThreads(int)"}
}
func UpdateThreadsFlag(threads *int) UpdateOption {
	return UpdateOption{func(opts *updateOptionImpl) {
		if threads == nil {
			return
		}
		opts.has_threads = true
		opts.threads = *threads
	}, "gen.UpdateThreads(int)"}
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
		opt.f(res)
	}
	return res
}

func MakeUpdateOptions(opts ...UpdateOption) UpdateOptions {
	return makeUpdateOptionImpl(opts...)
}
