package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/andreyvit/diff"
)

func TestGenOpts(t *testing.T) {
	var tests = []struct {
		name       string
		optType    string
		implType   string
		fieldSpecs []string
		opts       []GenOptsOption
	}{
		{
			name:    "empty",
			optType: "SomeOption",
		},
		{
			name:     "impl",
			optType:  "SomeOption",
			implType: "explicitImpl",
		},
		{
			name:       "fields",
			optType:    "SomeOption",
			fieldSpecs: []string{"foo", "bar:string", "baz:float64"},
		},
		{
			name:       "prefix",
			optType:    "SomeOption",
			fieldSpecs: []string{"foo", "bar:string", "baz:float64"},
			opts: []GenOptsOption{
				GenOptsPrefix("Prefix"),
			},
		},
		{
			name:       "prefixOptsType",
			optType:    "SomeOption",
			fieldSpecs: []string{"foo", "bar:string", "baz:float64"},
			opts: []GenOptsOption{
				GenOptsPrefixOptsType(true),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			goldenFile := path.Join("testdata", "goldens", fmt.Sprintf("%s.go.golden", test.name))
			goldenBytes, err := ioutil.ReadFile(goldenFile)
			if err != nil {
				t.Fatalf("reading file: %s: %v", goldenFile, err)
			}
			golden := string(goldenBytes)

			home, err := os.UserHomeDir()
			if err != nil {
				t.Fatalf("finding home: %v", err)
			}
			goImportsBin := path.Join(home, "go", "bin", "goimports")

			got, err := GenOpts(test.optType, test.implType, ".", goImportsBin, test.fieldSpecs, test.opts...)
			if err != nil {
				t.Fatalf("GenOpts(%q,%q,%v): %v", test.optType, test.implType, test.fieldSpecs, err)
			}
			gotWithPackage := "package goldens\n" + got
			tmp := goldenFile + ".tmp"
			defer os.Remove(tmp)
			gotFormatted, err := gofmt(gotWithPackage, tmp)
			if err != nil {
				t.Fatalf("formatting %s: %v", gotWithPackage, err)
			}

			// set to true to update goldens
			if false {
				if err := ioutil.WriteFile(goldenFile, []byte(gotFormatted), 0755); err != nil {
					t.Fatalf("updating golden %s: %v", goldenFile, err)
				}
			}

			if want, got := golden, gotFormatted; want != got {
				fmt.Println(got)
				t.Errorf("GenOpts(%q,%q,%v) want != got:\n\n------\n%s\n-------", test.optType, test.implType, test.fieldSpecs, diff.LineDiff(want, got))
			}
		})
	}
}

func gofmt(contents, tmp string) (string, error) {
	if err := ioutil.WriteFile(tmp, []byte(contents), 0755); err != nil {
		return "", err
	}
	cmd := exec.Command("gofmt", tmp)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	res := buf.String()
	return res, nil
}
