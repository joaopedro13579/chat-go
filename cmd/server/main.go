package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Message struct {
	chatid int
	text   string
}
type User struct {
	ID          int
	chatid      int
	Name        string
	Address     string
	destination int
	Conn        net.Conn
	chats       []*Chat
}
type Chat struct {
	ID       int
	Name     string
	users    []*User
	messages []Message
}

var (
	users  = []*User{}
	chats  = []*Chat{}
	mutex  sync.Mutex
	nextID = 1
)

func login(name string, conn net.Conn) *User {
	mutex.Lock()
	defer mutex.Unlock()

	for _, user := range users {
		if user.Name == name {
			user.Conn = conn
			user.Address = conn.RemoteAddr().String()
			user.destination = 0
			return user
		}
	}

	user := &User{
		ID:      len(users) + len(chats) + 1,
		Name:    name,
		Address: conn.RemoteAddr().String(),
		Conn:    conn,
	}
	fmt.Printf("%s connected from %s\n", user.Name, user.Address)
	users = append(users, user)
	return user
}
func finduserbyid(id int) *User {
	for _, user := range users {
		if user.ID == id {
			return user
		}
	}
	return nil
}
func finduser(name string) *User {
	for _, user := range users {
		if user.Name == name {
			return user
		}
	}
	return nil
}
func findchat(id int) *Chat {
	for _, chat := range chats {
		if chat.ID == id {
			return chat
		}
	}
	return nil
}
func createChat(name string, user *User) *Chat {
	mutex.Lock()
	var temp []*User
	temp = append(temp, user)
	defer mutex.Unlock()
	chat := &Chat{
		ID:    len(users) + len(chats) + 1,
		Name:  name,
		users: temp,
	}
	chats = append(chats, chat)
	return chat
}
func findreciever(id int) *[]*User {
	for _, user := range users {
		if user.ID == id {
			temp := []*User{}
			temp = append(temp, user)
			return &temp
		}
	}
	for _, chat := range chats {
		if chat.ID == id {
			return &chat.users
		}
	}
	return nil
}
func findchatbyname(name string) *Chat {
	for _, chat := range chats {
		if chat.Name == name {
			return chat
		}
	}
	return nil
}
func broadcast(sender *User, message Message) {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Printf("senderinfo: %s\n destination: %d", sender.Name, sender.destination)
	text := fmt.Sprintf("[%s]: %s\n", sender.Name, message.text)
	if sender.destination != 0 {
		for _, user := range *findreciever(sender.destination) {
			if user.Name != sender.Name {
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
	for i, u := range users {
		if u == user {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	mutex.Unlock()

	fmt.Printf("%s disconnected\n", user.Name)
	broadcast(user, Message{chatid: 0, text: "left the chat"})
}
func findChatByUserID(userID int) *Chat {
	for _, chat := range chats {
		for _, user := range chat.users {
			if user.ID == userID {
				return chat
			}
		}
	}
	return nil
}
func message_processer(message string, user *User) *Message {
	fmt.Print(message)
	if strings.HasPrefix(message, "/to ") {
		target := finduser(strings.TrimPrefix(message, "/to "))
		if target == nil {
			user.Conn.Write([]byte("User not found\n"))
			return nil
		}

		user.destination = target.ID
		return nil
	} else if strings.HasPrefix(message, "/tochat ") {
		user.destination = findchatbyname(strings.TrimPrefix(message, "/tochat ")).ID
		chat := findchatbyname(strings.TrimPrefix(message, "/tochat "))
		if chat == nil {
			user.Conn.Write([]byte("Chat not found\n"))
			return nil
		}

		user.destination = chat.ID
		return nil
	} else if strings.HasPrefix(message, "/whereami") {
		if user.destination == 0 {
			user.Conn.Write([]byte("You are in the main chat\n"))
			return nil
		} else if findchat(user.destination) != nil {
			user.Conn.Write([]byte(fmt.Sprintf("You are in chat %s (ID: %d)\n", findchat(user.destination).Name, user.destination)))
			return nil
		} else if finduserbyid(user.destination) != nil {
			user.Conn.Write([]byte(fmt.Sprintf("You are chatting with %s (ID: %d)\n", finduserbyid(user.destination).Name, user.destination)))
			return nil
		}

	} else if strings.HasPrefix(message, "/createChat") {
		chat := createChat(strings.TrimPrefix(message, "/createChat "), user)
		user.destination = chat.ID
		user.Conn.Write([]byte(fmt.Sprintf("Chat %s created with ID %d\n", chat.Name, chat.ID)))
		return nil
	} else if strings.HasPrefix(message, "/addToChat") {
		params := strings.Fields(strings.TrimPrefix(message, "/addToChat "))

		if len(params) != 2 {
			user.Conn.Write([]byte(
				"Usage: /addToChat <username> <chatname>\n",
			))
			return nil
		}

		addedUser := finduser(params[0])
		if addedUser == nil {
			user.Conn.Write([]byte("User not found\n"))
			return nil
		}

		chat := findchatbyname(params[1])
		if chat == nil {
			user.Conn.Write([]byte("Chat not found\n"))
			return nil
		}

		chat.users = append(chat.users, addedUser)

		if addedUser.Conn != nil {
			addedUser.Conn.Write([]byte(fmt.Sprintf(
				"You have been added to chat %s (ID: %d)\n",
				chat.Name,
				chat.ID,
			)))
		}

		user.Conn.Write([]byte(fmt.Sprintf(
			"User %s added to chat %s\n",
			addedUser.Name,
			chat.Name,
		)))

		return nil
	} else if strings.HasPrefix(message, "/listchats") {
		temp := chats
		for _, chat := range temp {
			for _, _user := range chat.users {
				if _user.ID == user.ID {
					user.Conn.Write([]byte(fmt.Sprintf("Chat: %s (ID: %d)\n", chat.Name, chat.ID)))
				}
			}
			return nil
		}
	} else if strings.HasPrefix(message, "/changeTo") {
		target := strings.TrimPrefix(message, "/changeTo ")

		if u := finduser(target); u != nil {
			user.destination = u.ID
			return nil
		}

		if c := findchatbyname(target); c != nil {
			user.destination = c.ID
			return nil
		}

		user.Conn.Write([]byte("Destination not found\n"))
		return nil
	} else if strings.HasPrefix(message, "/help") {
		help := `
				Available commands:

				/help
					Show this help message

				/to <username>
					Start a private conversation with a user

				/tochat <chatname>
					Send messages to a chat room

				/whereami
					Shows your current destination

				/createChat <chatname>
					Create a new chat room and switch to it

				/addToChat <username> <chatname>
					Add a user to a chat room

				/listchats
					List all chats you belong to

				/changeTo <username|chatname>
					Change your destination to a user or chat

				Typing any normal message sends it to the current destination.
				If no destination is selected, messages go to the main chat.
				`
		user.Conn.Write([]byte(help))
		return nil
	} else {
		message := &Message{
			chatid: user.destination,
			text:   message,
		}
		fmt.Printf("Processed message: %s\n", message.text)
		return message
		// ... existing command handling ...
	}
	return nil
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
	message := Message{
		chatid: 0,
		text:   fmt.Sprintf("%s has joined the chat", user.Name),
	}
	broadcast(user, message)

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
		_message := message_processer(message, user)
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
