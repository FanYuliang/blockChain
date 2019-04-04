package endpoints

import (
	"encoding/json"
	"mp2/blockchain"
	"mp2/node_membership"
)

type Endpoint struct {
	Types     [] string
	FEndpoint FailureDetectionMeta
	TEndpoint TransactionMeta
	BEndpoint BlockMeta
}

type FailureDetectionMeta struct {
	Type             int // 0: join, 1: ping, 2: ack 3: leave
	Record           []node_membership.Entry
	InitialTimeStamp int64
	IpAddress        string
	//Transactions     map[string]blockchain.Transaction
	//Block 			 blockchain.Block

}

type TransactionMeta struct {
	Tx  	[]blockchain.Transaction
}

type BlockMeta struct {
	Block  	blockchain.Block
}

func (e *Endpoint)  ToBytes() []byte {
	res, _ := json.Marshal(e)
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