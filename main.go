package main

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spudtrooper/genopts/genopts"
	"github.com/spudtrooper/genopts/options"
)

var (
	optType        = flag.String("opt_type", "Option", "The name of the primary options type")
	implType       = flag.String("impl_type", "", "The name of the implementation type; if empty this is derived from --opts_type")
	prefixOptsType = flag.Bool("prefix_opts_type", false, "Prefix each option function with the --opts_type; --prefix takes precendence over --prefix_opts_type")
	prefix         = flag.String("prefix", "", "Prefix each option with this string; --prefix takes precendence over --prefix_opts_type")
	outfile        = flag.String("outfile", "", "Output result to this file in addition to printing to STDOUT")
)

func genOpts() error {
	if *optType == "" {
		return errors.Errorf("--opt_type required")
	}
	output, err := genopts.GenOpts(*optType, *implType, flag.Args(),
		options.Prefix(*prefix),
		options.PrefixOptsType(*prefixOptsType),
		options.Outfile(*outfile))
	if err != nil {
		return err
	}
	if output != "" {
		fmt.Println(output)
	}
	return nil
}

func main() {
	flag.Parse()
	if err := genOpts(); err != nil {
		panic(err)
	}
}
