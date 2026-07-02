package manager

var Help = `
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
