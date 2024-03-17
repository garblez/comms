package message

import (
	"comms/user"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Message struct {
	Content string    `json:"content"`
	From    user.User `json:"from"`
}

func NewMessage(content string, from user.User) Message {
	return Message{
		Content: content,
		From:    from,
	}
}

func (m Message) String() string {
	return fmt.Sprintf("[%s] %s", m.From.Alias, m.Content)
}

func Transmit(message Message, recipient user.User) {
	conn, err := net.DialTCP("tcp", nil, &recipient.Address)
	if err != nil { // <-- refuses to connect at this point here!
		log.Println(err.Error())
		return
	}
	defer conn.Close()

	bytes, _ := json.Marshal(message) // What could seriously go wrong in marshalling a known struct?
	n, err := conn.Write(bytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("%s (%d bytes)\n", string(bytes[:]), n)
}
