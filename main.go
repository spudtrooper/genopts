package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spudtrooper/genopts/genopts"
	"github.com/spudtrooper/genopts/gitversion"
	"github.com/spudtrooper/genopts/options"
	"github.com/spudtrooper/goutil/check"
	"github.com/spudtrooper/goutil/io"
	"github.com/spudtrooper/goutil/or"
)

var (
	optType        = flag.String("opt_type", "", "The name of the primary options type. If empty and there is a prefix, --opt_type is the prefix + \"Option\"; if prefix is empty --opt_type is \"Option\"")
	implType       = flag.String("impl_type", "", "The name of the implementation type; if empty this is derived from --opts_type")
	prefixOptsType = flag.Bool("prefix_opts_type", false, "Prefix each option function with the --opts_type; --prefix takes precendence over --prefix_opts_type")
	prefix         = flag.String("prefix", "", "Prefix each option with this string; --prefix takes precendence over --prefix_opts_type")
	outfile        = flag.String("outfile", "", "Output result to this file in addition to printing to STDOUT")
	update         = flag.Bool("update", false, "update all files recurisvely in the current directory or directory specified by --update_dir")
	updateDir      = flag.String("update_dir", ".", "directory for update")
	updateFile     = flag.String("update_file", "", "single file to update")
	goimports      = flag.String("goimports", "", "full path to goimports, if empty we use ~/go/bin/goimports")
	excludeDirs    = flag.String("exclude_dirs", "", "comma-separated list of directory base names to exclude when --update is set")
	config         = flag.String("config", "", "absolute location of config. If empty we'll look in $update_dir/.genopts")
	writeConfig    = flag.Bool("write_config", false, "update the expected config file. This is used to set the config after setting explicit flags")
	quiet          = flag.Bool("quiet", false, "quite logging")
)

type Config struct {
	ExcludedDirs []string
	GoImports    string
}

func (c Config) Empty() bool {
	return len(c.ExcludedDirs) == 0 && c.GoImports == ""
}

func findConfig() (Config, error) {
	var configFile string
	if f := path.Join(*updateDir, ".genopts"); io.FileExists(f) {
		configFile = f
	} else if io.FileExists(*config) {
		configFile = *config
	}
	if configFile == "" {
		return Config{}, nil
	}
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return Config{}, err
	}
	if !*quiet {
		log.Printf("using config from %s", configFile)
	}
	return config, nil
}

func realMain() error {
	if gitversion.CheckVersionFlag() {
		return nil
	}
	cfg, err := findConfig()
	if err != nil {
		return err
	}

	goImportsBin := or.String(*goimports, cfg.GoImports)
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

	if *writeConfig {
		c := Config{
			ExcludedDirs: strings.Split(*excludeDirs, ","),
			GoImports:    *goimports,
		}
		bytes, err := json.Marshal(c)
		if err != nil {
			return err
		}
		configFile := path.Join(*updateDir, ".genopts")
		if err := ioutil.WriteFile(configFile, bytes, 0755); err != nil {
			return err
		}
		log.Printf("wrote config to %s", configFile)
	} else if !*quiet && cfg.Empty() && (*goimports != "" || *excludeDirs != "") {
		expectedConfigFile := or.String(*config, path.Join(*updateDir, ".genopts"))
		fmt.Printf("***\n")
		fmt.Printf("*** To run with this configuration without passing explicit flags,\n")
		fmt.Printf("*** run the same command adding the --write_config flag or copy the\n")
		fmt.Printf("*** following to %s\n", expectedConfigFile)
		fmt.Printf("***\n")
	}

	if *update || len(os.Args) == 1 {
		bin, err := os.Executable()
		if err != nil {
			return err
		}
		excludedDirs := cfg.ExcludedDirs
		if *excludeDirs != "" {
			for _, dir := range strings.Split(*excludeDirs, ",") {
				dir = strings.TrimSpace(dir)
				excludedDirs = append(excludedDirs, dir)
			}
		}
		if err := genopts.UpdateDir(dir, bin, goImportsBin, excludedDirs); err != nil {
			return err
		}
		return nil
	}

	if *updateFile != "" {
		bin, err := os.Executable()
		if err != nil {
			return err
		}
		if err := genopts.UpdateFile(*updateFile, bin, goImportsBin); err != nil {
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
