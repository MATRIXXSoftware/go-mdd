package dictionary

import (
	"testing"

	"github.com/matrixxsoftware/go-mdd/mdd/field"
	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	matrixxSchema := &Configuration{
		Containers: []Container{
			{
				ID:                   "Container1",
				Key:                  10001,
				CreatedSchemaVersion: 5200,
				DeletedSchemaVersion: 5270,
				Fields: []Field{
					{
						ID:                   "Field1",
						Datatype:             "string",
						CreatedSchemaVersion: 5200,
						DeletedSchemaVersion: 5270,
					},
					{
						ID:                   "Field2",
						Datatype:             "bool",
						CreatedSchemaVersion: 5200,
						DeletedSchemaVersion: 5260,
					},
					{
						ID:                   "Field3",
						Datatype:             "int32",
						CreatedSchemaVersion: 5200,
						DeletedSchemaVersion: 5270,
						IsArray:              true,
					},
				},
			},
		},
	}

	dict := NewWithSchema(matrixxSchema, nil)

	// Not Found, SchemaVersion not in range
	def, found := dict.Lookup(10001, 5280, 1)
	assert.False(t, found)

	// Not Found, Key not err
	def, found = dict.Lookup(10002, 5262, 1)
	assert.False(t, found)

	// Found
	def, found = dict.Lookup(10001, 5262, 1)
	assert.True(t, found)
	assert.Equal(t, 10001, def.Key)
	assert.Equal(t, 5262, def.SchemaVersion)
	assert.Equal(t, 1, def.ExtVersion)
	assert.Equal(t, "Container1", def.Name)
	assert.Equal(t, 2, len(def.Fields))

	// Field1
	assert.Equal(t, 0, def.Fields[0].Number)
	assert.Equal(t, "Field1", def.Fields[0].Name)
	assert.Equal(t, field.String, def.Fields[0].Type)
	assert.False(t, def.Fields[0].IsMulti)
	assert.False(t, def.Fields[0].IsContainer)

	// Field2 is out of version range

	// Field3
	assert.Equal(t, 1, def.Fields[1].Number)
	assert.Equal(t, "Field3", def.Fields[1].Name)
	assert.Equal(t, field.Int32, def.Fields[1].Type)
	assert.True(t, def.Fields[1].IsMulti)
	assert.False(t, def.Fields[1].IsContainer)
}

func TestLookup2(t *testing.T) {
	config := &Configuration{
		Containers: []Container{
			{
				ID:                   "Container2",
				Key:                  10002,
				CreatedSchemaVersion: 5260,
				Fields: []Field{
					{
						ID:       "Field1",
						Datatype: "string",
					},
					{
						ID:       "Field2",
						Datatype: "bool",
					},
				},
			},
		},
	}

	dict := NewWithSchema(config, nil)

	// Not Found, SchemaVersion not in range
	def, found := dict.Lookup(10002, 5250, 1)
	assert.False(t, found)

	// Found
	def, found = dict.Lookup(10002, 5262, 1)
	assert.True(t, found)
	assert.Equal(t, 10002, def.Key)
	assert.Equal(t, 5262, def.SchemaVersion)
	assert.Equal(t, 1, def.ExtVersion)
	assert.Equal(t, "Container2", def.Name)
	assert.Equal(t, 2, len(def.Fields))

	// Field1
	assert.Equal(t, 0, def.Fields[0].Number)
	assert.Equal(t, "Field1", def.Fields[0].Name)
	assert.Equal(t, field.String, def.Fields[0].Type)
	assert.False(t, def.Fields[0].IsMulti)
	assert.False(t, def.Fields[0].IsContainer)

	// Field2
	assert.Equal(t, 1, def.Fields[1].Number)
	assert.Equal(t, "Field2", def.Fields[1].Name)
	assert.Equal(t, field.Bool, def.Fields[1].Type)
	assert.False(t, def.Fields[1].IsMulti)
	assert.False(t, def.Fields[1].IsContainer)
}
