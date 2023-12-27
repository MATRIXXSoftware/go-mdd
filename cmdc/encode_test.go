package cmdc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

func TestEncode(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("1")},
					{Data: []byte("20")},
					{Data: []byte("(5:three)")},
					{Data: []byte("400000")},
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,20,(5:three),400000]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncode2(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("1")},
					{Data: []byte("20")},
					{Data: []byte("(5:three)")},
					{Data: []byte("4000")},
				},
			},
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    20,
					Depth:         0,
					Key:           -7,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("")},
					{Data: []byte("200")},
					{Data: []byte("(7:FooBar3)")},
					{Data: []byte("")},
					{Data: []byte("100000")},
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,20,(5:three),4000]<1,20,0,-7,5222,2>[,200,(7:FooBar3),,100000]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncodeEmptyField(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncodeKnownType(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Type: field.Int32, Value: int32(1)},
					{Type: field.String, Value: "three"},
					{Type: field.String, Value: "富爸"},
					{Data: []byte("4000")},
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,(5:three),(6:富爸),4000]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncodeNested(t *testing.T) {

	subContainers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    20,
					Depth:         0,
					Key:           -7,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					{Data: []byte("")},
					{Data: []byte("200")},
					{Data: []byte("(7:FooBar3)")},
					{Data: []byte("")},
					{Data: []byte("100000")},
				},
			},
		},
	}

	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{
					Version:       1,
					TotalField:    18,
					Depth:         0,
					Key:           -6,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []mdd.Field{
					*mdd.NewBasicField(int8(1)),
					*mdd.NewStructField(codec, &subContainers),
					*mdd.NewBasicField("three"),
					{Data: []byte("(4:four)")},
					*mdd.NewBasicField(uint32(5000)),
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,<1,20,0,-7,5222,2>[,200,(7:FooBar3),,100000],(5:three),(4:four),5000]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestEncodedNested2(t *testing.T) {
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{Version: 1, TotalField: 5, Depth: 0, Key: 200, SchemaVersion: 5222, ExtVersion: 2},
				Fields: []mdd.Field{
					*mdd.NewNullField(field.Int32),
					*mdd.NewBasicField(int32(100)),
					*mdd.NewStructField(codec,
						&mdd.Containers{
							Containers: []mdd.Container{
								{
									Header: mdd.Header{Version: 1, TotalField: 3, Depth: 0, Key: 100, SchemaVersion: 5222, ExtVersion: 2},
									Fields: []mdd.Field{
										*mdd.NewBasicField(uint32(0)),
										*mdd.NewBasicField("OK"),
										*mdd.NewNullField(field.UInt32),
									},
								},
							},
						},
					),
					*mdd.NewBasicField(uint64(2000000000)),
				},
			},
		},
	}

	expected := "<1,5,0,200,5222,2>[,100,<1,3,0,100,5222,2>[0,(2:OK),],2000000000]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}
