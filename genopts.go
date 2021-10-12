package main

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spudtrooper/genopts/genopts"
)

var (
	optsType = flag.String("opts_type", "", "The name of the primary options type")
	implType = flag.String("impl_type", "", "The name of the implementation type; if empty this is derived from --opts_type")
)

func genOpts() error {
	if *optsType == "" {
		return errors.Errorf("--opts_type required")
	}
	output, err := genopts.GenOpts(*optsType, *implType, flag.Args())
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func main() {
	flag.Parse()
	if err := genOpts(); err != nil {
		panic(err)
	}
}
