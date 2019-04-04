package endpoints

import (
	"encoding/json"
	"log"
	"mp2/blockchain"
	"mp2/node_membership"
	"os"
)

type Endpoint struct {
	Types     [] int
	FEndpoint FailureDetectionMeta
	TEndpoint TransactionMeta
	BEndpoint BlockMeta
	REndpoint RequestMissingTransactionMeta
}

type RequestMissingTransactionMeta struct{
	Type 			int // 0 receive, 1 reply
	RequesterIPaddr string
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

func (e *Endpoint)SetEndpointType(endpointTypes ...string) {
	e.Types = make([] int, 1)
	for _, endpointType := range endpointTypes {
		if endpointType == "FailureDetection" {
			e.Types = append(e.Types, 0)
		} else if endpointType == "Transaction" {
			e.Types = append(e.Types, 1)
		} else if endpointType == "Block" {
			e.Types = append(e.Types, 2)
		} else if endpointType == "RequestMissingTransaction"{
			e.Types = append(e.Types,3)
		} else {
			log.Fatal("Bad endpoint type!")
			os.Exit(12)
		}
	}
}

func (e *Endpoint)GetEndpointTypes() []string {
	res := make([] string, len(e.Types))
	for _, endpointType := range e.Types {
		if endpointType == 0 {
			res = append(res, "FailureDetection")
		} else if endpointType == 1 {
			res = append(res, "Transaction")
		} else if endpointType == 2 {
			res = append(res, "Block")
		} else if endpointType == 3{
			res = append(res, "RequestMissingTransaction")
		}
	}
	return  res
}

func EncodeFailureDetectionActionType(endpointType string) int {
	if endpointType == "Join" {
		return 0
	} else if endpointType == "Ping" {
		return 1
	} else if endpointType == "Ack" {
		return 2
	} else if endpointType == "Quit" {
		return 3
	}
	return -1
}