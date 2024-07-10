package dictionary

import (
	"encoding/xml"
	"io"
)

type Configuration struct {
	Name     xml.Name  `xml:"configuration"`
	Subtypes []Subtype `xml:"subtype"`
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

func Parse(data []byte) (*Configuration, error) {
	var config Configuration
	err := xml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func Load(reader io.Reader) (*Dictionary, error) {
	return nil, nil
}
