package server

import (
	"encoding/json"
	"mp2/blockchain"
)

type Endpoint struct {
	EndpointType     int // 0: join, 1: ping, 2: ack 3: leave
	Record           []Entry
	InitialTimeStamp int64
	IpAddress        string
	//Transactions     map[string]blockchain.Transaction
	Block 			 blockchain.Block
}

func (a *Endpoint)  ToBytes() []byte {
	res, _ := json.Marshal(a)
	return res
}

func EncodeEndpointType(endpointType string) int {
	if endpointType == "Join" {
		return 0
	} else if endpointType == "Ping" {
		return 1
	} else if endpointType == "Ack" {
		return 2
	} else if endpointType == "QUIT" {
		return 3
	} else if endpointType == "New Block" {
		return 4
	}
	return -1
}