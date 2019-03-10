package server

import (
	"bufio"
	"fmt"
	"net"
)

type Server struct {
	Name              string
	MyAddress         string
}

func (s *Server) TalkWithServiceServer(serviceConn net.Conn) {
	for {
		message, _ := bufio.NewReader(serviceConn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message), "\n")
	}
}