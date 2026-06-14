package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	message, _ := bufio.NewReader(conn).ReadString('\n')

	fmt.Println("Mensagem do servidor:")
	fmt.Println(message)
}
