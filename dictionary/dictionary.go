package dictionary

import "github.com/matrixxsoftware/go-mdd/mdd/field"

type ContainerDefinition struct {
	Name          string
	Key           int
	SchemaVersion int
	ExtVersion    int
	Fields        []FieldDefinition
}

type FieldDefinition struct {
	Number      int
	Name        string
	Type        field.Type
	IsMulti     bool
	IsContainer bool
}

type compositeKey struct {
	key           int
	schemaVersion int
	extVersion    int
}

type Dictionary struct {
	definitions map[compositeKey]ContainerDefinition
}

func New() *Dictionary {
	return &Dictionary{
		definitions: make(map[compositeKey]ContainerDefinition),
	}
}

func (d *Dictionary) Lookup(key, schemaVersion, extVersion int) (*ContainerDefinition, bool) {
	ckey := compositeKey{
		key:           key,
		schemaVersion: schemaVersion,
		extVersion:    extVersion,
	}
	return d.get(ckey)
}

// Add RWLock in future
func (d *Dictionary) get(ckey compositeKey) (*ContainerDefinition, bool) {
	containerDefinition, ok := d.definitions[ckey]
	return &containerDefinition, ok
}

func (d *Dictionary) Add(def *ContainerDefinition) {
	ckey := compositeKey{
		key:           def.Key,
		schemaVersion: def.SchemaVersion,
		extVersion:    def.ExtVersion,
	}
	d.definitions[ckey] = *def
}
