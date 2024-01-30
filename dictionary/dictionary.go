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

type Dictionary struct {
	definitions map[int]ContainerDefinition
}

func New() *Dictionary {
	return &Dictionary{
		definitions: make(map[int]ContainerDefinition),
	}
}

// Add RWLock in future
func (d *Dictionary) Get(key int) (*ContainerDefinition, bool) {
	containerDefinition, ok := d.definitions[key]
	return &containerDefinition, ok
}

func (d *Dictionary) Add(containerDefinition *ContainerDefinition) {
	d.definitions[containerDefinition.Key] = *containerDefinition
}
