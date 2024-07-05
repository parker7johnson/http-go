package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

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

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	request, _ := createRequest(string(buffer))

	log.Println(request)
	if request.Path == "/" {

		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Split(request.Path, "/")[1] == "echo" {

		message := strings.Split(request.Path, "/")[2]
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)))
	} else if strings.Split(request.Path, "/")[1] == "user-agent" {

		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(request.Headers["User-Agent"]), request.Headers["User-Agent"])))
	} else if strings.Contains(request.Path, "file") {
		fileDir := os.Args[2]
		bytes, err := os.ReadFile(fmt.Sprintf("%s%s", fileDir, strings.Split(request.Path, "/")[2]))
		if err == nil {
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(bytes), string(bytes))))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

	} else {

		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func createRequest(buffer string) (*HTTPRequest, error) {
	splitRequest := strings.Split(buffer, "\r\n")
	request := &HTTPRequest{}
	request.Method = strings.Split(splitRequest[0], " ")[0]
	request.Path = strings.Split(splitRequest[0], " ")[1]
	request.Body = splitRequest[len(splitRequest)-1]
	headers := make(map[string]string)
	for i := 0; i < len(splitRequest); i++ {
		if !strings.Contains(splitRequest[i], ":") {
			continue
		}

		splitHeader := strings.Split(splitRequest[i], " ")
		headers[strings.TrimSuffix(splitHeader[0], ":")] = splitHeader[1]

	}
	request.Headers = headers

	return request, nil
}
