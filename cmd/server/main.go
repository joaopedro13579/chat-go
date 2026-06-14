package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type User struct {
	ID          int
	Name        string
	Address     string
	destination int
	Conn        net.Conn
}

var (
	users  = make(map[string]*User)
	mutex  sync.Mutex
	nextID = 1
)

func login(name string, conn net.Conn) *User {
	mutex.Lock()
	defer mutex.Unlock()

	if user, exists := users[name]; exists {
		user.Conn = conn
		user.Address = conn.RemoteAddr().String()
		user.destination = 0
		return user
	}

	user := &User{
		ID:      nextID,
		Name:    name,
		Address: conn.RemoteAddr().String(),
		Conn:    conn,
	}

	users[name] = user
	nextID++

	fmt.Printf("%s connected from %s\n", user.Name, user.Address)

	return user
}
func finduser(name string) *User {
	for _, user := range users {
		if user.Name == name {
			return user
		}
	}
	user := &User{
		ID:      0,
		Name:    "",
		Address: "",
		Conn:    nil,
	}
	return user
}
func broadcast(sender *User, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	text := fmt.Sprintf("[%s]: %s\n", sender.Name, message)
	if sender.destination != 0 {
		for _, user := range users {
			if user.Name != sender.Name && user.ID == sender.destination {
				_, _ = user.Conn.Write([]byte(text))
			}
		}
	} else {
		for _, user := range users {
			if user.Name != sender.Name {
				_, _ = user.Conn.Write([]byte(text))
			}
		}
	}
}
func removeUser(user *User) {
	mutex.Lock()
	delete(users, user.Name)
	mutex.Unlock()

	fmt.Printf("%s disconnected\n", user.Name)
	broadcast(user, "left the chat")
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

	user := login(name, conn)

	_, _ = conn.Write([]byte(
		fmt.Sprintf("Welcome %s! Your ID is %d\n", user.Name, user.ID),
	))

	broadcast(user, "joined the chat")

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			removeUser(user)
			return
		}

		message = strings.TrimSpace(message)

		if message == "" {
			continue
		}
		if strings.HasPrefix(message, "/to ") {
			if strings.TrimPrefix(message, "/to ") == "" {
				user.destination = 0
			} else {
				user.destination = finduser(strings.TrimPrefix(message, "/to ")).ID
				fmt.Println(user.destination)
			}
		}
		broadcast(user, message)
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
