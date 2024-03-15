package main

import (
	"comms/chat"
	"comms/message"
	"comms/user"
	"fmt"
	"time"
)

func main() {
	alice := user.NewUser("Alice", "", 8080)
	bob := user.NewUser("Bob", "", 8888)
	chat := chat.NewChat("Test Chat", bob)
	defer chat.Close()
	msg := message.NewMessage("Watson, get in here I need you!", alice)

	chat.Open()

	fmt.Println(chat.Title)
	for i := 0; i < 10; i++ {
		<-time.After(time.Second)
		chat.Send(msg)
	}

	<-time.After(time.Millisecond)
	chat.Print()
}
