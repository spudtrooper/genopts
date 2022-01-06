package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spudtrooper/genopts/genopts"
	"github.com/spudtrooper/genopts/options"
	"github.com/spudtrooper/goutil/check"
)

var (
	optType        = flag.String("opt_type", "Option", "The name of the primary options type")
	implType       = flag.String("impl_type", "", "The name of the implementation type; if empty this is derived from --opts_type")
	prefixOptsType = flag.Bool("prefix_opts_type", false, "Prefix each option function with the --opts_type; --prefix takes precendence over --prefix_opts_type")
	prefix         = flag.String("prefix", "", "Prefix each option with this string; --prefix takes precendence over --prefix_opts_type")
	outfile        = flag.String("outfile", "", "Output result to this file in addition to printing to STDOUT")
	update         = flag.Bool("update", false, "update all files recurisvely in the current directory or directory specified by --update_dir")
	updateDir      = flag.String("update_dir", ".", "directory for update")
	goimports      = flag.String("goimports", "", "full path to goimports, if empty we use ~/go/bin/goimports")
)

func realMain() error {
	goImportsBin := *goimports
	if goImportsBin == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		goImportsBin = path.Join(home, "go", "bin", "goimports")
	}
	dir, err := filepath.Abs(*updateDir)
	if err != nil {
		return err
	}

	if *update {
		bin, err := os.Executable()
		if err != nil {
			return err
		}
		if err := genopts.UpdateDir(dir, bin, goImportsBin); err != nil {
			return err
		}
		return nil
	}
	if err := genOpts(dir, goImportsBin); err != nil {
		return err
	}
	return nil
}

func genOpts(dir, goImportsBin string) error {
	if *optType == "" {
		return errors.Errorf("--opt_type required")
	}
	output, err := genopts.GenOpts(*optType, *implType, dir, goImportsBin, flag.Args(),
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
	check.Err(realMain())
}
