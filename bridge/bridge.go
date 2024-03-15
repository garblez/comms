package bridge

import (
	"log"
	"net"
)

type Bridge struct {
	Listener net.Listener
}

func NewBridge() (Bridge, error) {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err.Error())
		return Bridge{}, err
	}

	return Bridge{
		Listener: l,
	}, nil
	// Deliberately leave out defer l.Close(); we want to handle this on the chat struct
}

func (b *Bridge) Close() {
	b.Listener.Close()
}
