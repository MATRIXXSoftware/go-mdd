package cmdc

import (
	"testing"
)

var codec = NewCodec()

func BenchmarkDecode(b *testing.B) {
	mdc := "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100],4]"
	for i := 0; i < b.N; i++ {
		codec.Decode([]byte(mdc))
	}
}

func BenchmarkEncode(b *testing.B) {
	mdc := "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100],4]"
	containers, _ := codec.Decode([]byte(mdc))
	for i := 0; i < b.N; i++ {
		codec.Encode(containers)
	}
}
