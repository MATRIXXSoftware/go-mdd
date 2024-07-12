package dictionary

import (
	"encoding/xml"
	"io"
)

type Configuration struct {
	Name       xml.Name    `xml:"configuration"`
	Subtypes   []Subtype   `xml:"subtype"`
	Containers []Container `xml:"container"`
}

type Subtype struct {
	ID       string  `xml:"id,attr"`
	Datatype string  `xml:"datatype"`
	Values   []Value `xml:"value"`
}

type Value struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type Container struct {
	ID                   string        `xml:"id,attr"`
	DocDescription       string        `xml:"doc_description"`
	Key                  int           `xml:"key"`
	CreatedSchemaVersion int           `xml:"created_schema_version"`
	DeletedSchemaVersion int           `xml:"deleted_schema_version"`
	BaseContainer        BaseContainer `xml:"base_container,omitempty"`
	Fields               []Field       `xml:"field"`
}

type BaseContainer struct {
	ID   string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

type Field struct {
	ID                   string `xml:"id,attr"`
	DocDescription       string `xml:"doc_description"`
	Datatype             string `xml:"datatype"`
	IsArray              bool   `xml:"array"`
	IsList               bool   `xml:"list"`
	StructID             string `xml:"struct_id"`
	SubTypeReference     string `xml:"subtype_reference"`
	CreatedSchemaVersion int    `xml:"created_schema_version"`
	DeletedSchemaVersion int    `xml:"deleted_schema_version"`
}

func Parse(reader io.Reader) (*Configuration, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = xml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
