package mdd

import (
	"fmt"

	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

type ClientTransport interface {
	Send([]byte) ([]byte, error)
	Close() error
}

type ServerTransport interface {
	Listen() error
	Handler(handler func([]byte) ([]byte, error))
	Close() error
}

type Client struct {
	Codec     Codec
	Transport ClientTransport
}

func (c *Client) SendMessage(request *Containers) (*Containers, error) {

	reqBody, err := c.Codec.Encode(request)
	if err != nil {
		return nil, err
	}

	respBody, err := c.Transport.Send(reqBody)
	if err != nil {
		return nil, err
	}

	response, err := c.Codec.Decode(respBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type Server struct {
	Codec     Codec
	Transport ServerTransport
}

func (s *Server) MessageHandler(handler func(*Containers) (*Containers, error)) {

	h := func(reqBody []byte) ([]byte, error) {
		request, err := s.Codec.Decode(reqBody)
		if err != nil {
			return nil, err
		}

		// hopId, err := extractHopId(request)
		// if err != nil {
		// 	return nil, err
		// }
		// fmt.Printf("HopId: %d\n", hopId)

		response, err := handler(request)
		if err != nil {
			return nil, err
		}

		respBody, err := s.Codec.Encode(response)
		if err != nil {
			return nil, err
		}

		return respBody, nil
	}

	s.Transport.Handler(h)
}

func extractHopId(containers *Containers) (uint32, error) {
	// Get MtxMsg container (key 93)
	mtxMsg := containers.GetContainer(93)
	if mtxMsg == nil {
		return 0, fmt.Errorf("container MtxMsg is missing")
	}

	// Assume no changes to the position of hopId field
	f := mtxMsg.GetField(14)

	if f.Data == nil {
		return 0, fmt.Errorf("hopId field is missing")
	}

	// Copy the field data to a new field
	hopIdField := Field{
		Data:  f.Data,
		Type:  field.UInt32,
		Codec: f.Codec,
	}

	// Get the value of the field
	hopId, err := hopIdField.GetValue()
	if err != nil {
		return 0, err
	}

	return hopId.(uint32), nil
}
