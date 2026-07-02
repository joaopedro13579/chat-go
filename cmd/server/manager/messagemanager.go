package manager

import (
	"fmt"
	"network-project/cmd/server/structs"
	"strings"
)

func Messageprocesser(message string, user *structs.User) *structs.Message {
	fmt.Print(message)
	if strings.HasPrefix(message, "/to ") {
		target := Finduser(strings.TrimPrefix(message, "/to "))
		if target == nil {
			user.Conn.Write([]byte("User not found\n"))
			return nil
		}

		user.Destination = target.ID
		return nil
	} else if strings.HasPrefix(message, "/tochat ") {
		user.Destination = Findchatbyname(strings.TrimPrefix(message, "/tochat ")).ID
		chat := Findchatbyname(strings.TrimPrefix(message, "/tochat "))
		if chat == nil {
			user.Conn.Write([]byte("Chat not found\n"))
			return nil
		}

		user.Destination = chat.ID
		return nil
	} else if strings.HasPrefix(message, "/whereami") {
		if user.Destination == 0 {
			user.Conn.Write([]byte("You are in the main chat\n"))
			return nil
		} else if Findchat(user.Destination) != nil {
			user.Conn.Write([]byte(fmt.Sprintf("You are in chat %s (ID: %d)\n", Findchat(user.Destination).Name, user.Destination)))
			return nil
		} else if Finduserbyid(user.Destination) != nil {
			user.Conn.Write([]byte(fmt.Sprintf("You are chatting with %s (ID: %d)\n", Finduserbyid(user.Destination).Name, user.Destination)))
			return nil
		}

	} else if strings.HasPrefix(message, "/createChat") {
		chat := CreateChat(strings.TrimPrefix(message, "/createChat "), user)
		user.Destination = chat.ID
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

		addedUser := Finduser(params[0])
		if addedUser == nil {
			user.Conn.Write([]byte("User not found\n"))
			return nil
		}

		chat := Findchatbyname(params[1])
		if chat == nil {
			user.Conn.Write([]byte("Chat not found\n"))
			return nil
		}

		chat.Users = append(chat.Users, addedUser)

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
		temp := Chats
		for _, chat := range temp {
			for _, _user := range chat.Users {
				if _user.ID == user.ID {
					user.Conn.Write([]byte(fmt.Sprintf("Chat: %s (ID: %d)\n", chat.Name, chat.ID)))
				}
			}
			return nil
		}
	} else if strings.HasPrefix(message, "/changeTo") {
		target := strings.TrimPrefix(message, "/changeTo ")

		if u := Finduser(target); u != nil {
			user.Destination = u.ID
			return nil
		}

		if c := Findchatbyname(target); c != nil {
			user.Destination = c.ID
			return nil
		}

		user.Conn.Write([]byte("Destination not found\n"))
		return nil
	} else if strings.HasPrefix(message, "/help") {
		user.Conn.Write([]byte(Help))
		return nil
	} else if strings.HasPrefix(message, "/get") {
		param := strings.Fields(message)
		url := param[0]
		url = url
	} else {
		message := &structs.Message{
			Chatid: user.Destination,
			Text:   message,
		}
		fmt.Printf("Processed message: %s\n", message.Text)
		return message
		// ... existing command handling ...
	}
	return nil
}
