package dictionary

import (
	"fmt"
	"sync"

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
	definitions     map[compositeKey]ContainerDefinition
	matrixxSchema   *Configuration
	extensionSchema *Configuration
	mu              sync.RWMutex
}

func New() *Dictionary {
	return &Dictionary{
		definitions:     make(map[compositeKey]ContainerDefinition),
		matrixxSchema:   nil,
		extensionSchema: nil,
	}
}

func NewWithSchema(matrixxSchema *Configuration, extensionSchema *Configuration) *Dictionary {
	return &Dictionary{
		definitions:     make(map[compositeKey]ContainerDefinition),
		matrixxSchema:   matrixxSchema,
		extensionSchema: extensionSchema,
	}
}

func stringToType(datatype string) (field.Type, error) {
	switch datatype {
	case "struct":
		return field.Struct, nil
	case "string":
		return field.String, nil
	case "bool":
		return field.Bool, nil
	case "signed int8":
		return field.Int8, nil
	case "signed int16":
		return field.Int16, nil
	case "signed int32":
		return field.Int32, nil
	case "signed int64":
		return field.Int64, nil
	case "signed int128":
		return field.Int128, nil
	case "unsigned int8":
		return field.UInt8, nil
	case "unsigned int16":
		return field.UInt16, nil
	case "unsigned int32":
		return field.UInt32, nil
	case "unsigned int64":
		return field.UInt64, nil
	case "unsigned int128":
		return field.UInt128, nil
	case "decimal":
		return field.Decimal, nil
	case "date":
		return field.Date, nil
	case "time":
		return field.Time, nil
	case "datetime":
		return field.DateTime, nil
	case "blob":
		return field.Blob, nil
	case "buffer id":
		return field.BufferID, nil
	case "field key":
		return field.FieldKey, nil
	case "phone number":
		return field.PhoneNo, nil
	case "object id":
		return field.ObjectID, nil
	default:
		return field.Unknown, fmt.Errorf("Unknown datatype: %s", datatype)
	}
}

func (d *Dictionary) search(key, schemaVersion, extVersion int) (*ContainerDefinition, error) {
	var container Container
	var isFound, isPrivate bool

	filterContainer := func(containers []Container, version int) bool {
		for _, c := range containers {
			if c.Key == key &&
				(c.CreatedSchemaVersion == 0 || version >= c.CreatedSchemaVersion) &&
				(c.DeletedSchemaVersion == 0 || version < c.DeletedSchemaVersion) {
				container = c
				isFound = true
				return true
			}
		}
		return false
	}

	if d.matrixxSchema != nil {
		isFound = filterContainer(d.matrixxSchema.Containers, schemaVersion)
	}

	if !isFound && d.extensionSchema != nil {
		isFound = filterContainer(d.extensionSchema.Containers, extVersion)
		isPrivate = isFound
	}

	if !isFound {
		return nil, fmt.Errorf("Container not found: key=%d, schemaVersion=%d, extVersion=%d",
			key, schemaVersion, extVersion)
	}

	filterField := func(fields []Field, version int) ([]FieldDefinition, error) {
		var result []FieldDefinition
		number := 0

		for _, f := range fields {
			dataType, err := stringToType(f.Datatype)
			if err != nil {
				return nil, fmt.Errorf("Error field %s Container key=%d, schemaVersion=%d, extVersion=%d: %v",
					f.ID, key, schemaVersion, extVersion, err)
			}
			if (f.CreatedSchemaVersion == 0 || version >= f.CreatedSchemaVersion) &&
				f.DeletedSchemaVersion == 0 || version < f.DeletedSchemaVersion {
				result = append(result, FieldDefinition{
					Number:      number,
					Name:        f.ID,
					Type:        dataType,
					IsMulti:     f.IsList || f.IsArray,
					IsContainer: f.StructID != "",
				})
				number++
			}
		}
		return result, nil
	}

	var fields []FieldDefinition
	var err error
	if isPrivate {
		fields, err = filterField(container.Fields, extVersion)
	} else {
		fields, err = filterField(container.Fields, schemaVersion)
	}

	if err != nil {
		return nil, err
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

func (d *Dictionary) Lookup(key, schemaVersion, extVersion int) (*ContainerDefinition, error) {

	ckey := compositeKey{
		key:           key,
		schemaVersion: schemaVersion,
		extVersion:    extVersion,
	}
	result, found := d.get(ckey)

	if !found {
		result, err := d.search(key, schemaVersion, extVersion)
		if err == nil {
			d.Add(result)
			return result, nil
		}
		return nil, err
	}

	return result, nil
}

func (d *Dictionary) get(ckey compositeKey) (*ContainerDefinition, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	containerDefinition, ok := d.definitions[ckey]
	return &containerDefinition, ok
}

func (d *Dictionary) Add(def *ContainerDefinition) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ckey := compositeKey{
		key:           def.Key,
		schemaVersion: def.SchemaVersion,
		extVersion:    def.ExtVersion,
	}
	d.definitions[ckey] = *def
}
