package cmdc

import "github.com/matrixxsoftware/go-mdd/mdd"

type Cmdc struct {
}

func NewCodec() *Cmdc {
	return &Cmdc{}
}

func (c *Cmdc) Decode(data []byte) (*mdd.Containers, error) {
	return Decode(data)
}

func (c *Cmdc) Encode(containers *mdd.Containers) ([]byte, error) {
	return Encode(containers)
}
