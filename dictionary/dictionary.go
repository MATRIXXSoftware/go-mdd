package dictionary

import (
	"fmt"

	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

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
	definitions   map[compositeKey]ContainerDefinition
	configuration *Configuration // Can be nil
}

func New() *Dictionary {
	return &Dictionary{
		definitions:   make(map[compositeKey]ContainerDefinition),
		configuration: nil,
	}
}

func NewWithConfig(config *Configuration) *Dictionary {
	return &Dictionary{
		definitions:   make(map[compositeKey]ContainerDefinition),
		configuration: config,
	}
}

func stringToType(datatype string) (field.Type, error) {
	switch datatype {
	case "string":
		return field.String, nil
	case "bool":
		return field.Bool, nil
	case "int8":
		return field.Int8, nil
	case "int32":
		return field.Int32, nil
		// TODO
	default:
		return field.Unknown, fmt.Errorf("Unknown datatype: %s", datatype)
	}
}

func (d *Dictionary) Search(key, schemaVersion, extVersion int) (*ContainerDefinition, error) {
	isFound := false
	var container Container

	// Construct new ContainerDefinition from Base Schema Configuration
	if d.configuration != nil {
		for _, c := range d.configuration.Containers {
			if c.Key == key &&
				(c.CreatedSchemaVersion == 0 || schemaVersion >= c.CreatedSchemaVersion) &&
				(c.DeletedSchemaVersion == 0 || schemaVersion < c.DeletedSchemaVersion) {
				container = c
				isFound = true
			}
		}
	}

	// TODO Construct new ContainerDefinition from Extension Schema

	fields := []FieldDefinition{}
	number := 0
	for _, f := range container.Fields {
		dataType, err := stringToType(f.Datatype)
		if err != nil {
			return nil, err
		}

		if (f.CreatedSchemaVersion == 0 || schemaVersion >= f.CreatedSchemaVersion) &&
			f.DeletedSchemaVersion == 0 || schemaVersion < f.DeletedSchemaVersion {
			fieldDefinition := FieldDefinition{
				Number:      number,
				Name:        f.ID,
				Type:        dataType,
				IsMulti:     f.IsList || f.IsArray,
				IsContainer: f.StructID != "",
			}
			fields = append(fields, fieldDefinition)
			number++
		}
	}

	if !isFound {
		return nil, fmt.Errorf("Container not found: key=%d, schemaVersion=%d, extVersion=%d",
			key, schemaVersion, extVersion)
	}

	def := &ContainerDefinition{
		Name:          container.ID,
		Key:           container.Key,
		SchemaVersion: schemaVersion,
		ExtVersion:    extVersion,
		Fields:        fields,
	}
	return def, nil
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
