package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
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
	Encoding    string
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
			Encoding:    req.Header.Get("Accept-Encoding"),
		}

		handlerResponse(c, 200, content)
	} else if strings.HasPrefix(path, "/files") && req.Method == "GET" {
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

	} else if strings.HasPrefix(path, "/files") && req.Method == "POST" {
		fileName := strings.TrimPrefix(path, "/files/")

		dir := os.Args[2]

		buf := make([]byte, req.ContentLength)

		n, _ := req.Body.Read(buf)

		data := string(buf[:n])

		os.WriteFile(dir+fileName, []byte(data), 0644)

		handlerResponseCreated(c)

	} else {
		handlerResponseNotFound(c)
	}
}

func handlerResponse(c net.Conn, statusCode int, content *Content) {

	statusReason := getReasonByStatusCode(statusCode)

	if content == nil {
		c.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", statusCode, statusReason)))
	} else if content.Encoding == "" || !strings.Contains(content.Encoding, "gzip") {
		c.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", statusCode, statusReason, content.ContentType, content.Length, content.Body)))
	} else {

		var b bytes.Buffer

		gz := gzip.NewWriter(&b)

		_, err := gz.Write([]byte(content.Body))

		if err != nil {
			log.Fatal(err)
		}

		if err := gz.Close(); err != nil {
			log.Fatal(err)
		}

		c.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Encoding: gzip\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", statusCode, statusReason, content.ContentType, len(b.String()), b.String())))
	}

}

func handlerResponseOk(c net.Conn) {
	handlerResponse(c, 200, nil)
}

func handlerResponseCreated(c net.Conn) {
	handlerResponse(c, 201, nil)
}

func handlerResponseNotFound(c net.Conn) {
	handlerResponse(c, 404, nil)
}

func getReasonByStatusCode(statusCode int) string {

	switch statusCode {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 404:
		return "Not Found"
	}

	return ""
}
