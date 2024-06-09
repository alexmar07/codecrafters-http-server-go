package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	conn, err := l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handler(conn)
	// conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

func handler(c net.Conn) {

	buffer := make([]byte, 1024)

	n, _ := c.Read(buffer)

	req := string(buffer[:n])

	splits := strings.Split(req, " ")

	if len(splits) < 3 {
		log.Fatal("Not excepted request")
	}

	path := splits[1]

	switch path {
	case "/":
		c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}
