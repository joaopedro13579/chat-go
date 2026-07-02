package structs

import (
	"net"
)

type User struct {
	ID          int
	Chatid      int
	Name        string
	Address     string
	Destination int
	Conn        net.Conn
	Chats       []*Chat
}
