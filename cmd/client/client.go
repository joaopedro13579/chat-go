package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	var serverAddr string
	
	if len(os.Args) > 1 {
		serverAddr = os.Args[1] + ":8080"
	} else {
		fmt.Print("Enter server IP (or press Enter for localhost): ")
		serverAddr, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		serverAddr = strings.TrimSpace(serverAddr)
		if serverAddr == "" {
			serverAddr = "127.0.0.1"
		}
		serverAddr = serverAddr + ":8080"
	}
	
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()


	fmt.Printf("Connected to server at %s\n", serverAddr)
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
