package genopts

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spudtrooper/genopts/log"
	"github.com/spudtrooper/goutil/io"
	"github.com/spudtrooper/goutil/sets"
)

var (
	updateDirBlacklist = map[string]bool{
		"// genopts {{.CommandLine}}": true,
	}
)

func UpdateDir(dir, bin, goImportsBin string, excludedDirs []string) error {
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

	for f, cmdLine := range filesAndCommandLines {
		if err := updateFile(f, bin, goImportsBin, dir, cmdLine); err != nil {
			return err
		}
	}

	return nil
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
	args = append(args, "--quiet")
	for _, arg := range strings.Split(cmdLine, " ") {
		arg = removeQuotes(arg)
		args = append(args, arg)
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
		if strings.HasPrefix(line, "// genopts") && !updateDirBlacklist[line] {
			cmdLine := strings.TrimSpace(strings.Replace(line, "// genopts", "", 1))
			return cmdLine, nil
		}
	}
	return "", nil
}
