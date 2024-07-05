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

	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	print(string(buffer))
	request, _ := createRequest(string(buffer))
	log.Println()
	log.Println(request)
	if request.Method == "GET" {
		handleGETRequest(request, conn)
	} else if request.Method == "POST" {
		handlePOSTRequest(request, conn)
	}

	conn.Close()
}

func handlePOSTRequest(request *HTTPRequest, conn net.Conn) {
	if strings.Contains(request.Path, "/files/") {
		writeFile(*request)

		conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))

	}
}

func handleGETRequest(request *HTTPRequest, conn net.Conn) {
	if request.Path == "/" {

		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Split(request.Path, "/")[1] == "echo" {

		message := strings.Split(request.Path, "/")[2]
		headersSlice := strings.Split(request.Headers["Accept-Encoding"], ", ")
		for i, v := range headersSlice {
			// v = strings.Trim(v, ",")
			println(i, v)
		}
		println("valid encoding index out at ")
		println(find(headersSlice, "gzip"))
		if find(headersSlice, "gzip") != -1 {
			print("entered writing content encoding header")
			encodingIndex := find(headersSlice, "gzip")
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: %s\r\n\r\n", headersSlice[encodingIndex])))
		} else {
			print("did not write content encoding header")
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)))

		}
	} else if strings.Split(request.Path, "/")[1] == "user-agent" {

		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(request.Headers["User-Agent"]), request.Headers["User-Agent"])))
	} else if strings.Contains(request.Path, "/files/") {
		bytes, err := readFile(request)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(bytes), string(bytes))))
		}

	} else {

		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func readFile(request *HTTPRequest) ([]byte, error) {
	fileDir := os.Args[2]
	return os.ReadFile(fmt.Sprintf("%s%s", fileDir, strings.Split(request.Path, "/")[2]))
}

func writeFile(request HTTPRequest) error {
	fileDir := os.Args[2]
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		os.Mkdir(fileDir, 0755)
	}
	outFile := fmt.Sprintf("%s%s", fileDir, strings.Split(request.Path, "/")[2])
	println(outFile)
	return os.WriteFile(outFile, []byte(request.Body), 0666)
}

func createRequest(buffer string) (*HTTPRequest, error) {
	splitRequest := strings.Split(buffer, "\r\n")
	splitRequest[len(splitRequest)-1] = strings.Trim(splitRequest[len(splitRequest)-1], "\x00")
	request := &HTTPRequest{}
	request.Method = strings.Split(splitRequest[0], " ")[0]
	request.Path = strings.Split(splitRequest[0], " ")[1]
	request.Body = splitRequest[len(splitRequest)-1]
	headers := make(map[string]string)
	for i := 0; i < len(splitRequest); i++ {
		if !strings.Contains(splitRequest[i], ":") {
			continue
		}

		if strings.Contains(splitRequest[i], "Accept-Encoding") {
			header := strings.Split(splitRequest[i], ":")
			headers[header[0]] = header[1]
		} else {
			splitHeader := strings.Split(splitRequest[i], " ")
			headers[strings.TrimSuffix(splitHeader[0], ":")] = splitHeader[1]
		}

	}
	request.Headers = headers

	return request, nil
}

func find(slice []string, value string) int {
	for i := range slice {

		if strings.Contains(slice[i], value) {
			return i
		}
	}
	return -1
}
