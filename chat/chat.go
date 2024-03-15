package chat

import (
	"comms/bridge"
	"comms/message"
	"comms/user"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Chat struct {
	Title    string
	history  []message.Message
	Members  []user.User
	pending  chan message.Message
	RealTime bool
	bridge   bridge.Bridge
}

func NewChat(title string, members ...user.User) Chat {
	b, err := bridge.NewBridge()
	if err != nil {
		log.Fatal("Failed to open bridge")
	}

	return Chat{
		Title:   title,
		history: make([]message.Message, 0),
		Members: members,
		pending: make(chan message.Message),
		bridge:  b,
	}
}

func (c *Chat) Close() {
	close(c.pending)
	c.bridge.Close()
}

func (c *Chat) Print() {
	fmt.Println("------- ", c.Title, " -------")
	for _, message := range c.history {
		fmt.Println(message)
	}
}

func (c *Chat) Open() {
	fmt.Println("Starting chat", c.Title)

	// Start monitoring the chat, adding any pending messages to its history until the chat is closed
	// This could also be where we split off to listen to other users - locally, we can just send a message
	// to the chat directly from a user
	go func() {
		for newMessage := range c.pending {
			// Keep adding any new, incoming messages to the chat log ad-infinitum until the chat
			c.history = append(c.history, newMessage)
			if c.RealTime {
				fmt.Println(newMessage)
			}
		}
	}()

	go c.receiveMessages()
	fmt.Println("Goroutine started for receiving messages")
}

func (chat *Chat) receive(c net.Conn) {
	log.Println("Connection opened for incoming message.")
	buf := make([]byte, 2048)
	n, err := c.Read(buf)
	if err != nil {
		log.Println("Failed to read from connection")
		return
	}

	data := message.Message{}
	err = json.Unmarshal(buf[:n], &data)
	if err != nil {
		log.Println("Error unmarshalling:", err.Error())
	}

	chat.pending <- data
	c.Close()
}

func (chat *Chat) receiveMessages() {
	for {
		conn, err := chat.bridge.Listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
			break
		}

		fmt.Println("Starting receiver")

		go chat.receive(conn)
	}
}

func (c *Chat) Send(msg message.Message) {
	for _, member := range c.Members {
		log.Println("Sending to", member.Alias)
		message.Transmit(msg, member)
	}

	c.pending <- msg
}
