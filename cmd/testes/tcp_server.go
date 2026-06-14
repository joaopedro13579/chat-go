package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}

	fmt.Println("Servidor TCP iniciado na porta 8080 ✔️")

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		fmt.Println("Novo cliente conectado:", conn.RemoteAddr())

		conn.Write([]byte("Conexao realizada com sucesso!\n"))

		conn.Close()
	}
}
