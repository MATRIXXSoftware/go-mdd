package transport

import (
	"testing"
	"time"
)

func TestTransport(t *testing.T) {
	go func() {
		err := MddServer("localhost:8080")
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	msg := ExampleMessage{Field1: 1, Field2: 2}

	err := MddClient("localhost:8080", &msg)
	if err != nil {
		panic(err)
	}

	time.Sleep(100 * time.Millisecond)
}
