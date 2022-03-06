package gen

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spudtrooper/genopts/log"
	goutilerrors "github.com/spudtrooper/goutil/errors"
	"github.com/spudtrooper/goutil/io"
	"github.com/spudtrooper/goutil/or"
	"github.com/spudtrooper/goutil/parallel"
	"github.com/spudtrooper/goutil/sets"
)

var (
	updateDirBlacklist = map[string]bool{
		"// genopts {{.CommandLine}}":                 true,
		"//go:" + "generate genopts {{.CommandLine}}": true,
	}
)

func UpdateDir(dir, bin, goImportsBin string, excludedDirs []string, uOpts ...UpdateOption) error {
	opts := MakeUpdateOptions(uOpts...)
	threads := or.Int(opts.Threads(), 10)

	excludedDirSet := sets.String(excludedDirs)
	filesAndCommandLines := map[string]string{}
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && excludedDirSet[filepath.Base(path)] {
			return filepath.SkipDir
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			cmdLine, err := extractCommandLine(path)
			if err != nil {
				return err
			}
			if cmdLine != "" {
				filesAndCommandLines[path] = cmdLine
			}
		}
		return nil
	}); err != nil {
		return err
	}

	files := make(chan interface{})
	go func() {
		var sorted []string
		for f := range filesAndCommandLines {
			sorted = append(sorted, f)
		}
		sort.Strings(sorted)
		for _, f := range sorted {
			files <- f
		}
		close(files)
	}()

	col := goutilerrors.MakeSyncErrorCollector()

	parallel.ExecAndDrain(files, threads, func(i interface{}) (interface{}, error) {
		f := i.(string)
		cmdLine := filesAndCommandLines[f]
		if err := updateFile(f, bin, goImportsBin, dir, cmdLine); err != nil {
			col.Add(err)
		}
		return nil, nil
	})

	return col.Build()
}

func UpdateFile(f, bin, goImportsBin string) error {
	if isGo := filepath.Ext(f) == ".go"; !isGo {
		log.Printf("%s is not a go file", f)
		return nil
	}
	cmdLine, err := extractCommandLine(f)
	if err != nil {
		return err
	}
	if cmdLine == "" {
		log.Printf("%s doesn't have an update commandline", f)
		return nil
	}

	if err := updateFile(f, bin, goImportsBin, ".", cmdLine); err != nil {
		return err
	}

	return nil
}

func updateFile(f, bin, goImportsBin, dir, cmdLine string) error {
	var args []string

	addOutfile := func() error {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.Errorf("os.Getwd: %v", err)
		}
		abs, err := filepath.Abs(f)
		if err != nil {
			return errors.Errorf("filepath.Abs(%q): %v", f, err)
		}
		rel, err := filepath.Rel(pwd, abs)
		if err != nil {
			return errors.Errorf("filepath.Rel(%q,%q): %v", pwd, abs, err)
		}
		args = append(args, "--outfile")
		args = append(args, rel)
		return nil
	}

	args = append(args, "--batch")
	if *verboseRun {
		args = append(args, "--verbose_run")
	}
	if err := addOutfile(); err != nil {
		return err
	}
	cmdLineParts := strings.Split(cmdLine, " ")
	for i := 0; i < len(cmdLineParts); {
		arg := cmdLineParts[i]
		i++
		arg = removeQuotes(arg)
		if strings.HasPrefix(arg, "--outfile=") {
			if err := addOutfile(); err != nil {
				return err
			}
		} else if arg == "--outfile" {
			i++
			if err := addOutfile(); err != nil {
				return err
			}
		} else {
			args = append(args, arg)
		}
	}

	if err := run(dir, bin, args...); err != nil {
		return err
	}
	if err := postGenCleanup(goImportsBin, dir, f); err != nil {
		return err
	}
	return nil
}

func removeQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	if len(s) > 0 && s[0] == '\'' && s[len(s)-1] == '\'' {
		s = s[1 : len(s)-1]
	}
	return s
}

func extractCommandLine(f string) (string, error) {
	lines, err := io.ReadLines(f)
	if err != nil {
		return "", err
	}
	for _, line := range lines {
		if (strings.HasPrefix(line, "// genopts") || strings.HasPrefix(line, "//go:"+"generate genopts")) && !updateDirBlacklist[line] {
			cmdLine := strings.TrimSpace(strings.Replace(line, "// genopts", "", 1))
			cmdLine = strings.TrimSpace(strings.Replace(cmdLine, "//go:"+"generate genopts", "", 1))
			return cmdLine, nil
		}
	}
	return "", nil
}
