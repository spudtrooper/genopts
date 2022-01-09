package genopts

import (
	"os"
	"path/filepath"
	"strings"

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
			cmdLine, err := checkForUpdate(path)
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
		var args []string
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

func checkForUpdate(f string) (string, error) {
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
