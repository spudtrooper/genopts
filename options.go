// DO NOT EDIT MANUALLY: Generated from https://github.com/spudtrooper/genopts
package genopts


//go:generate genopts --required  --outfile=options.go "gen/genopts.go" "gen/genopts_test.go" "gen/genoptsoptions.go" "gen/postgen.go" "gen/run.go" "gen/update.go" "gen/updateoptions.go"




type Option func(*optionImpl)

type Options interface {
	
Gen/genopts.go() bool
HasGen/genopts.go() bool	
Gen/genopts_test.go() bool
HasGen/genopts_test.go() bool	
Gen/genoptsoptions.go() bool
HasGen/genoptsoptions.go() bool	
Gen/postgen.go() bool
HasGen/postgen.go() bool	
Gen/run.go() bool
HasGen/run.go() bool	
Gen/update.go() bool
HasGen/update.go() bool	
Gen/updateoptions.go() bool
HasGen/updateoptions.go() bool
}

func Gen/genopts.go(gen/genopts.go bool) Option {
	return func(opts *optionImpl) {
		opts.has_gen/genopts.go = true
		opts.gen/genopts.go = gen/genopts.go
	}
}
func Gen/genopts.goFlag(gen/genopts.go *bool) Option {
	return func(opts *optionImpl) {
		if gen/genopts.go == nil {
			return
		}
		opts.has_gen/genopts.go = true
		opts.gen/genopts.go = *gen/genopts.go
	}
}

func Gen/genopts_test.go(gen/genopts_test.go bool) Option {
	return func(opts *optionImpl) {
		opts.has_gen/genopts_test.go = true
		opts.gen/genopts_test.go = gen/genopts_test.go
	}
}
func Gen/genopts_test.goFlag(gen/genopts_test.go *bool) Option {
	return func(opts *optionImpl) {
		if gen/genopts_test.go == nil {
			return
		}
		opts.has_gen/genopts_test.go = true
		opts.gen/genopts_test.go = *gen/genopts_test.go
	}
}

func Gen/genoptsoptions.go(gen/genoptsoptions.go bool) Option {
	return func(opts *optionImpl) {
		opts.has_gen/genoptsoptions.go = true
		opts.gen/genoptsoptions.go = gen/genoptsoptions.go
	}
}
func Gen/genoptsoptions.goFlag(gen/genoptsoptions.go *bool) Option {
	return func(opts *optionImpl) {
		if gen/genoptsoptions.go == nil {
			return
		}
		opts.has_gen/genoptsoptions.go = true
		opts.gen/genoptsoptions.go = *gen/genoptsoptions.go
	}
}

func Gen/postgen.go(gen/postgen.go bool) Option {
	return func(opts *optionImpl) {
		opts.has_gen/postgen.go = true
		opts.gen/postgen.go = gen/postgen.go
	}
}
func Gen/postgen.goFlag(gen/postgen.go *bool) Option {
	return func(opts *optionImpl) {
		if gen/postgen.go == nil {
			return
		}
		opts.has_gen/postgen.go = true
		opts.gen/postgen.go = *gen/postgen.go
	}
}

func Gen/run.go(gen/run.go bool) Option {
	return func(opts *optionImpl) {
		opts.has_gen/run.go = true
		opts.gen/run.go = gen/run.go
	}
}
func Gen/run.goFlag(gen/run.go *bool) Option {
	return func(opts *optionImpl) {
		if gen/run.go == nil {
			return
		}
		opts.has_gen/run.go = true
		opts.gen/run.go = *gen/run.go
	}
}

func Gen/update.go(gen/update.go bool) Option {
	return func(opts *optionImpl) {
		opts.has_gen/update.go = true
		opts.gen/update.go = gen/update.go
	}
}
func Gen/update.goFlag(gen/update.go *bool) Option {
	return func(opts *optionImpl) {
		if gen/update.go == nil {
			return
		}
		opts.has_gen/update.go = true
		opts.gen/update.go = *gen/update.go
	}
}

func Gen/updateoptions.go(gen/updateoptions.go bool) Option {
	return func(opts *optionImpl) {
		opts.has_gen/updateoptions.go = true
		opts.gen/updateoptions.go = gen/updateoptions.go
	}
}
func Gen/updateoptions.goFlag(gen/updateoptions.go *bool) Option {
	return func(opts *optionImpl) {
		if gen/updateoptions.go == nil {
			return
		}
		opts.has_gen/updateoptions.go = true
		opts.gen/updateoptions.go = *gen/updateoptions.go
	}
}

type optionImpl struct {
	gen/genopts.go bool
has_gen/genopts.go bool
	gen/genopts_test.go bool
has_gen/genopts_test.go bool
	gen/genoptsoptions.go bool
has_gen/genoptsoptions.go bool
	gen/postgen.go bool
has_gen/postgen.go bool
	gen/run.go bool
has_gen/run.go bool
	gen/update.go bool
has_gen/update.go bool
	gen/updateoptions.go bool
has_gen/updateoptions.go bool

}

func (o *optionImpl) Gen/genopts.go() bool { return o.gen/genopts.go }
func (o *optionImpl) HasGen/genopts.go() bool { return o.has_gen/genopts.go }
func (o *optionImpl) Gen/genopts_test.go() bool { return o.gen/genopts_test.go }
func (o *optionImpl) HasGen/genopts_test.go() bool { return o.has_gen/genopts_test.go }
func (o *optionImpl) Gen/genoptsoptions.go() bool { return o.gen/genoptsoptions.go }
func (o *optionImpl) HasGen/genoptsoptions.go() bool { return o.has_gen/genoptsoptions.go }
func (o *optionImpl) Gen/postgen.go() bool { return o.gen/postgen.go }
func (o *optionImpl) HasGen/postgen.go() bool { return o.has_gen/postgen.go }
func (o *optionImpl) Gen/run.go() bool { return o.gen/run.go }
func (o *optionImpl) HasGen/run.go() bool { return o.has_gen/run.go }
func (o *optionImpl) Gen/update.go() bool { return o.gen/update.go }
func (o *optionImpl) HasGen/update.go() bool { return o.has_gen/update.go }
func (o *optionImpl) Gen/updateoptions.go() bool { return o.gen/updateoptions.go }
func (o *optionImpl) HasGen/updateoptions.go() bool { return o.has_gen/updateoptions.go }



func makeOptionImpl(opts ...Option) *optionImpl {
	res := &optionImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeOptions(opts ...Option) Options {
	return makeOptionImpl(opts...)
}