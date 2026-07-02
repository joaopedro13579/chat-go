package manager

import (
	"network-project/cmd/server/structs"
)

func CreateChat(name string, user *structs.User) *structs.Chat {
	Mutex.Lock()
	var temp []*structs.User
	temp = append(temp, user)
	defer Mutex.Unlock()
	chat := &structs.Chat{
		ID:    len(Users) + len(Chats) + 1,
		Name:  name,
		Users: temp,
	}
	Chats = append(Chats, chat)
	return chat
}
