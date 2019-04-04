package server

import (
	"mp2/endpoints"
	"mp2/utils"
	"net"
	"strconv"
	"strings"
)

func (s *Server) sendMessageWithUDP(endpoint endpoints.Endpoint, ipAddress string) {
	//fmt.Println("ipAddress: ", ipAddress)
	arr := strings.Split(ipAddress, ":")

	myPort, err := strconv.Atoi(arr[1])
	utils.CheckError(err)

	iparr := utils.StringAddrToIntArr(ipAddress)
	Conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: iparr, Port: myPort, Zone: ""})
	utils.CheckError(err)
	defer Conn.Close()


	//fmt.Println("endpoint: ", endpoint)
	n, err := Conn.Write(endpoint.ToBytes())
	s.BandwidthLock.Lock()
	s.Bandwidth += float64(int(n)/1024)
	s.BandwidthLock.Unlock()
	utils.CheckError(err)
}



