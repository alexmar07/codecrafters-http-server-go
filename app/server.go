package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type Content struct {
	Length      int
	Body        string
	ContentType string
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {

		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handler(conn)
	}

}

func handler(c net.Conn) {

	req, err := http.ReadRequest(bufio.NewReader(c))

	if err != nil {
		log.Fatal(err)
	}

	path := req.URL.Path

	if path == "/" {
		handlerResponseOk(c)
	} else if path == "/user-agent" {

		userAgent := req.UserAgent()

		content := &Content{
			Length:      len(userAgent),
			Body:        userAgent,
			ContentType: "text/plain",
		}

		handlerResponse(c, 200, content)
	} else if strings.HasPrefix(path, "/echo") {
		param := strings.TrimPrefix(path, "/echo/")
		content := &Content{
			Length:      len(param),
			Body:        param,
			ContentType: "text/plain",
		}

		handlerResponse(c, 200, content)
	} else if strings.HasPrefix(path, "/files") {
		fileName := strings.TrimPrefix(path, "/files/")

		dir := os.Args[2]

		fileData, err := os.ReadFile(dir + fileName)

		if err != nil {
			handlerResponseNotFound(c)
		}

		content := &Content{
			Length:      len(fileData),
			Body:        string(fileData),
			ContentType: "application/octet-stream",
		}

		handlerResponse(c, 200, content)

	} else {
		handlerResponseNotFound(c)
	}
}

func handlerResponse(c net.Conn, statusCode int, content *Content) {

	statusReason := getReasonByStatusCode(statusCode)

	if content == nil {
		c.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", statusCode, statusReason)))
	} else {
		c.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", statusCode, statusReason, content.ContentType, content.Length, content.Body)))
	}

}

func handlerResponseOk(c net.Conn) {
	handlerResponse(c, 200, nil)
}

func handlerResponseNotFound(c net.Conn) {
	handlerResponse(c, 404, nil)
}

func getReasonByStatusCode(statusCode int) string {

	switch statusCode {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	}

	return ""
}
