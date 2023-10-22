package transport

import "net"

func MddClient(addr string, msg *ExampleMessage) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = msg.Encode(conn)
	if err != nil {
		return err
	}

	// Receive a response
	err = msg.Decode(conn)
	if err != nil {
		return err
	}
	return nil
}
