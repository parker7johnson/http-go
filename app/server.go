package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var c = make(chan int)

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

	opts := strings.Split(string(buffer), "\r\n")
	path := strings.Split(opts[0], " ")

	pathParts := strings.Split(path[1], "/")
	for i, v := range pathParts {
		println(i, v)
	}
	if path[1] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	} else if pathParts[1] == "echo" {
		message := pathParts[2]
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
	}
}
