package manager

import (
	"network-project/cmd/server/structs"
	"sync"
)

var (
	Users  = []*structs.User{}
	Chats  = []*structs.Chat{}
	Mutex  sync.Mutex
	NextID = 1
)

func Finduserbyid(id int) *structs.User {
	for _, user := range Users {
		if user.ID == id {
			return user
		}
	}
	return nil
}
func Finduser(name string) *structs.User {
	for _, user := range Users {
		if user.Name == name {
			return user
		}
	}
	return nil
}
func Findchat(id int) *structs.Chat {
	for _, chat := range Chats {
		if chat.ID == id {
			return chat
		}
	}
	return nil
}

func FindChatByUserID(userID int) *structs.Chat {
	for _, chat := range Chats {
		for _, user := range chat.Users {
			if user.ID == userID {
				return chat
			}
		}
	}
	return nil
}
func Findreciever(id int) *[]*structs.User {
	for _, user := range Users {
		if user.ID == id {
			temp := []*structs.User{}
			temp = append(temp, user)
			return &temp
		}
	}
	for _, chat := range Chats {
		if chat.ID == id {
			return &chat.Users
		}
	}
	return nil
}
func Findchatbyname(name string) *structs.Chat {
	for _, chat := range Chats {
		if chat.Name == name {
			return chat
		}
	}
	return nil
}
