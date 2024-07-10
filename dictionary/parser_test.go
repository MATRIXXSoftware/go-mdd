package dictionary

import (
	"testing"
)

func TestLoad(t *testing.T) {
	data := `
    <configuration>
        <subtype id='Enum5GDeliveryReportRequested'>
            <datatype>unsigned int16</datatype>
            <value id='1'>YES</value>
            <value id='2'>NO</value>
        </subtype>
    </configuration>
    `

	config, err := Parse([]byte(data))
	if err != nil {
		t.Error("Error parsing XML:", err)
	}
	t.Log("Parsed configuration:", config)
}
