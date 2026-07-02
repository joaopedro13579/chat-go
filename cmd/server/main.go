package main

import (
	"bufio"
	"fmt"
	"net"
	"network-project/cmd/server/log"
	"network-project/cmd/server/manager"
	"network-project/cmd/server/structs"
	"strings"
)

func broadcast(sender *structs.User, message structs.Message) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	log.WriteLog(fmt.Sprintf("[%s]: %s", sender.Name, message.Text))
	text := fmt.Sprintf("[%s]: %s\n", sender.Name, message.Text)
	if sender.Destination != 0 {
		for _, user := range *manager.Findreciever(sender.Destination) {
			if user.Name != sender.Name {
				_, _ = user.Conn.Write([]byte(text))
			}
		}
	} else {
		for _, user := range manager.Users {
			if user.Name != sender.Name {
				_, _ = user.Conn.Write([]byte(text))
			}
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	_, _ = conn.Write([]byte("Enter your name: "))

	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	name = strings.TrimSpace(name)

	if name == "" {
		_, _ = conn.Write([]byte("Invalid name\n"))
		return
	}

	user := manager.Login(name, conn)

	_, _ = conn.Write([]byte(
		fmt.Sprintf("Welcome %s! Your ID is %d\n", user.Name, user.ID),
	))
	message := structs.Message{
		Chatid: 0,
		Text:   fmt.Sprintf("%s has joined the chat", user.Name),
	}
	broadcast(user, message)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			manager.RemoveUser(user)
			broadcast(user, structs.Message{Chatid: 0, Text: "left the chat"})

			return
		}

		message = strings.TrimSpace(message)

		if message == "" {
			continue
		}
		_message := manager.Messageprocesser(message, user)
		if _message != nil {
			broadcast(user, *_message)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Chat server running on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}
