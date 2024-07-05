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

var statusCodes = make(map[string]string)

func main() {
	statusCodes = initMap()
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
	res := ""
	log.Println(request)
	if request.Path == "/" {
		res = request.createResponse(*request, "200")
		conn.Write([]byte(res))
	} else if strings.Split(request.Path, "/")[1] == "echo" {
		res = request.createResponse(*request, "200")
		conn.Write([]byte(res))
	} else if strings.Split(request.Path, "/")[1] == "user-agent" {
		res = request.createResponse(*request, "200")
		conn.Write([]byte(res))
	} else {
		res = request.createResponse(*request, "404")
		conn.Write([]byte(res))
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

func (*HTTPRequest) createResponse(req HTTPRequest, status string) string {
	res := ""
	headers := req.Headers
	if headers == nil {
		res = fmt.Sprintf("HTTP/1.1 %s %s\r\n\r\n", status, statusCodes[status])
	} else {
		res = fmt.Sprintf("HTTP/1.1 %s %s\r\n", status, statusCodes[status])
		for k, v := range headers {
			res += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		if strings.Contains(req.Path, "user-agent") {
			res += req.Headers["User-Agent"]
		} else if req.Body != "null" || req.Body != "" {
			res += req.Body
		}

		log.Default().Println(res)

	}
	return res
}

func initMap() map[string]string {
	statusCodes := map[string]string{
		"100": "Continue",
		"101": "Switching Protocols",
		"102": "Processing",
		"103": "Early Hints",
		"200": "OK",
		"201": "Created",
		"202": "Accepted",
		"203": "Non-Authoritative Information",
		"204": "No Content",
		"205": "Reset Content",
		"206": "Partial Content",
		"207": "Multi-Status",
		"208": "Already Reported",
		"226": "IM Used",
		"300": "Multiple Choices",
		"301": "Moved Permanently",
		"302": "Found",
		"303": "See Other",
		"304": "Not Modified",
		"305": "Use Proxy",
		"307": "Temporary Redirect",
		"308": "Permanent Redirect",
		"400": "Bad Request",
		"401": "Unauthorized",
		"402": "Payment Required",
		"403": "Forbidden",
		"404": "Not Found",
		"405": "Method Not Allowed",
		"406": "Not Acceptable",
		"407": "Proxy Authentication Required",
		"408": "Request Timeout",
		"409": "Conflict",
		"410": "Gone",
		"411": "Length Required",
		"412": "Precondition Failed",
		"413": "Payload Too Large",
		"414": "URI Too Long",
		"415": "Unsupported Media Type",
		"416": "Range Not Satisfiable",
		"417": "Expectation Failed",
		"418": "I'm a teapot",
		"421": "Misdirected Request",
		"422": "Unprocessable Entity",
		"423": "Locked",
		"424": "Failed Dependency",
		"425": "Too Early",
		"426": "Upgrade Required",
		"428": "Precondition Required",
		"429": "Too Many Requests",
		"431": "Request Header Fields Too Large",
		"451": "Unavailable For Legal Reasons",
		"500": "Internal Server Error",
		"501": "Not Implemented",
		"502": "Bad Gateway",
		"503": "Service Unavailable",
		"504": "Gateway Timeout",
		"505": "HTTP Version Not Supported",
		"506": "Variant Also Negotiates",
		"507": "Insufficient Storage",
		"508": "Loop Detected",
		"510": "Not Extended",
		"511": "Network Authentication Required",
	}

	return statusCodes
}
