package gen

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/spudtrooper/goutil/io"
)

const metadataFileName = "genopts-metadata.json"
const metadataDebug = true

type metadataField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type metadataType struct {
	RelativeFile string          `json:"relative_file"`
	Package      string          `json:"package"`
	TypeName     string          `json:"type_name"`
	Fields       []metadataField `json:"fields"`
}

type metadata struct {
	Types []metadataType `json:"types"`
}

type metadataCache struct {
	pwd string
	md  *metadata
	mu  sync.Mutex
}

func newMetadataCacheFromPWD(pwd string) *metadataCache {
	return &metadataCache{pwd: pwd}
}

func (m *metadataCache) init() error {
	if m.md != nil {
		return nil
	}
	m.md = &metadata{}
	m.mu.Lock()
	defer m.mu.Unlock()
	f := path.Join(m.pwd, metadataFileName)
	if !io.FileExists(f) {
		if metadataDebug {
			log.Printf("metadata doesn't exist: %s", f)
		}
		return nil
	}
	if metadataDebug {
		log.Printf("reading metadata from %s", f)
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	var res metadata
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	m.md = &res
	return nil
}

func (m *metadataCache) ForType(pkg, name string) (metadataType, error) {
	if err := m.init(); err != nil {
		return metadataType{}, err
	}
	for _, t := range m.md.Types {
		if t.Package == pkg && t.TypeName == name {
			return t, nil
		}
	}
	return metadataType{}, errors.Errorf("type %s.%s not found in metadata", pkg, name)
}

func (m *metadataCache) Update(t metadataType) error {
	if err := m.init(); err != nil {
		return err
	}
	m.md.Types = append(m.md.Types, t)
	b, err := json.MarshalIndent(m.md, "", "  ")
	if err != nil {
		return err
	}
	f := path.Join(m.pwd, metadataFileName)
	if metadataDebug {
		log.Printf("saving metadata to %s", f)
	}
	if err := ioutil.WriteFile(f, b, 0644); err != nil {
		return err
	}
	return nil
}
