package cmdc

import (
	"testing"
)

func testData() string {
	return sampleData2()
}

func BenchmarkDecode(b *testing.B) {
	data := testData()
	for i := 0; i < b.N; i++ {
		codec.Decode([]byte(data))
	}
}

func BenchmarkEncode(b *testing.B) {
	data := testData()
	containers, _ := codec.Decode([]byte(data))
	for i := 0; i < b.N; i++ {
		codec.Encode(containers)
	}
}
