package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")
	fmt.Println("Enter your name: ")
	myusername, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	conn.Write([]byte(myusername + "\n"))
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nDisconnected from server")
				os.Exit(0)
			}
			fmt.Print(message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, myusername) != true {
			_, err = conn.Write([]byte(text + "\n"))
		}
		if err != nil {
			fmt.Println("Send error:", err)
			return
		}
	}
}
