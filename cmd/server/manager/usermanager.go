package manager

import (
	"fmt"
	"net"
	"network-project/cmd/server/structs"
)

func Login(name string, conn net.Conn) *structs.User {
	Mutex.Lock()
	defer Mutex.Unlock()

	for _, user := range Users {
		if user.Name == name {
			user.Conn = conn
			user.Address = conn.RemoteAddr().String()
			user.Destination = 0
			return user
		}
	}

	user := &structs.User{
		ID:      len(Users) + len(Chats) + 1,
		Name:    name,
		Address: conn.RemoteAddr().String(),
		Conn:    conn,
	}
	fmt.Printf("%s connected from %s\n", user.Name, user.Address)
	Users = append(Users, user)
	return user
}
func RemoveUser(user *structs.User) {
	Mutex.Lock()
	for i, u := range Users {
		if u == user {
			Users = append(Users[:i], Users[i+1:]...)
			break
		}
	}
	Mutex.Unlock()

	fmt.Printf("%s disconnected\n", user.Name)
}
