package mdd

import (
	"fmt"
	"testing"
)

func TestMdd(t *testing.T) {
	// mdc := "<1,18,0,-6,5222,2>[,(6:value2),3,2021-09-07T08:00:25.000001Z," +
	// 	"2021-10-31,09:13:02.667997Z,88,5.5,"
	mdc := "<1,18,0,-6,5222>[1,abc,foo,bar]"
	encode()
	container := decode(mdc)
	fmt.Println("Decoded Container:")
	fmt.Printf("Header: %+v\n", container.Header)
	fmt.Printf("Fields: %+v\n", container.Fields)
}
