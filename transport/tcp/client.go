package tcp

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/matrixxsoftware/go-mdd/mdd"
	log "github.com/sirupsen/logrus"
)

type ClientTransport struct {
	conn       net.Conn
	Codec      mdd.Codec
	msgCache   map[uint32]chan *mdd.Containers
	msgMutex   sync.Mutex
	writeMutex sync.Mutex
	closeCh    chan struct{}
	closeWg    sync.WaitGroup
	seqId      uint32
}

func (c *ClientTransport) Close() error {
	close(c.closeCh)
	c.closeWg.Wait()
	return c.conn.Close()
}

func NewClientTransport(addr string, codec mdd.Codec) (*ClientTransport, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c := &ClientTransport{
		conn:     conn,
		Codec:    codec,
		msgCache: make(map[uint32]chan *mdd.Containers),
		closeCh:  make(chan struct{}),
		seqId:    0,
	}

	c.closeWg.Add(1)
	go func() {
		defer c.closeWg.Done()
		for {
			select {
			case <-c.closeCh:
				return
			default:
				respBody, err := Read(c.conn)
				if err != nil {
					if err == io.EOF {
						log.Infof("%s Connection closed", connStr(conn))
					} else if err == io.ErrUnexpectedEOF {
						log.Errorf("%s Connection closed unexpectedly", connStr(conn))
					} else {
						log.Errorf("%s Error reading from connection: %s", connStr(conn), err)
					}
					return
				}

				if respBody == nil {
					continue
				}

				// Handle the response
				err = c.processResponse(respBody)
				if err != nil {
					log.Errorf("%s %s", connStr(conn), err)
					return
				}
			}
		}
	}()

	return c, nil
}

func (c *ClientTransport) processResponse(respBody []byte) error {
	// Decode the response
	response, err := c.Codec.Decode(respBody)
	if err != nil {
		return err
	}

	hopId, err := extractHopId(response)
	if err != nil {
		return err
	}

	// Find the corresponding request
	c.msgMutex.Lock()
	ch, exists := c.msgCache[hopId]
	if exists {
		delete(c.msgCache, hopId)
	}
	c.msgMutex.Unlock()

	if !exists {
		return errors.New("unexpected response hopId")
	}

	// Send the response to the waiting channel
	ch <- response
	close(ch)

	return nil
}

func (c *ClientTransport) SendMessage(ctx context.Context, request *mdd.Containers) (*mdd.Containers, error) {

	hopId := atomic.AddUint32(&c.seqId, 1)

	err := injectHopId(request, hopId)
	if err != nil {
		return nil, err
	}

	reqBody, err := c.Codec.Encode(request)
	if err != nil {
		return nil, err
	}

	ch := make(chan *mdd.Containers, 1)

	c.msgMutex.Lock()
	c.msgCache[hopId] = ch
	c.msgMutex.Unlock()

	c.writeMutex.Lock()
	err = Write(c.conn, reqBody)
	c.writeMutex.Unlock()

	if err != nil {
		return nil, err
	}

	// Wait for the response from channel
	select {
	case response := <-ch:
		return response, nil
	case <-ctx.Done():
		c.msgMutex.Lock()
		delete(c.msgCache, hopId)
		c.msgMutex.Unlock()
		return nil, errors.New("timeout waiting for response")
	}
}
