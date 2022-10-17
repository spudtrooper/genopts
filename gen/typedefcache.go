package gen

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	//go : generate genopts --function RestaurantDetails --params verbose debugFailure
	genOptsFnRE  = regexp.MustCompile(`^//go.generate genopts (.*)`)
	extendsExtRE = regexp.MustCompile(`--extends[= ](\S+).*`)
)

type typeDefCache struct {
	dir   string
	types map[string]typeDef
}

func newTypeDefCache(dir string) *typeDefCache {
	return &typeDefCache{
		dir:   dir,
		types: map[string]typeDef{},
	}
}

func (t *typeDefCache) init() error {
	var goFiles []string
	files, err := ioutil.ReadDir(t.dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if ext := filepath.Ext(f.Name()); ext == ".go" {
			goFiles = append(goFiles, f.Name())
		}
	}

	var cmdLines []string
	for _, f := range goFiles {
		c, err := ioutil.ReadFile(filepath.Join(t.dir, f))
		if err != nil {
			return err
		}
		for _, line := range strings.Split(string(c), "\n") {
			if m := genOptsFnRE.FindStringSubmatch(line); len(m) == 2 {
				cmdLine := m[1]
				cmdLines = append(cmdLines, cmdLine)
			}
		}
	}

	for _, cmdLine := range cmdLines {
		r := csv.NewReader(strings.NewReader(cmdLine))
		r.Comma = ' '
		args, err := r.Read()
		if err != nil {
			return err
		}
		var extends []string
		if m := extendsExtRE.FindStringSubmatch(cmdLine); len(m) == 2 {
			extends = strings.Split(m[1], ",")
		}
		rest, name := findRest(args)
		fields := makeFields(rest)
		td := typeDef{
			name:    name,
			extends: extends,
			fields:  fields,
			args:    args,
		}
		t.types[name] = td
	}
	return nil
}

func (t *typeDefCache) findType(name string) (typeDef, error) {
	res, ok := t.types[name]
	if !ok {
		return typeDef{}, errors.Errorf("type %q not found", name)
	}
	return res, nil
}

func (t *typeDefCache) findExtendedTypes(extendsNames []string) ([]typeDef, error) {
	var res []typeDef
	for _, ext := range extendsNames {
		td, err := t.findType(ext)
		if err != nil {
			return nil, err
		}
		res = append(res, td)
		if err := exec.Command("genopts", td.args...).Run(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (tc *typeDefCache) transitiveFields(td typeDef) ([]string, error) {
	m := map[string]string{}
	if err := tc.transitiveFieldsLoop(td, m); err != nil {
		return nil, err
	}

	var res []string
	for f := range m {
		res = append(res, f)
	}
	return res, nil
}

func (tc *typeDefCache) transitiveFieldsLoop(td typeDef, res map[string]string) error {
	for _, f := range td.fields {
		n := title(f.Name)
		typ, ok := res[n]
		if ok {
			if typ != f.Type {
				return fmt.Errorf("field %q has conflicting types: %q and %q", f.Name, typ, f.Type)
			}
			continue
		}
		res[n] = f.Type
	}
	for _, e := range td.extends {
		t, err := tc.findType(e)
		if err != nil {
			return err
		}
		if err := tc.transitiveFieldsLoop(t, res); err != nil {
			return err
		}
	}
	return nil
}
