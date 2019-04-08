package endpoints

import (
	"encoding/json"
	"fmt"
	"os"
)

type Endpoint struct {
	Types      [] int
	FEndpoint  FailureDetectionMeta
	TEndpoint  TransactionMeta
	BEndpoint  BlockMeta
	RMEndpoint RequestMissingBlockMeta
}

func (e *Endpoint)  ToBytes() []byte {
	res, _ := json.Marshal(e)
	return res
}

func (e *Endpoint)SetEndpointType(endpointTypes ...string) {
	e.Types = make([] int, 1)
	for _, endpointType := range endpointTypes {
		if endpointType == "FailureDetection" {
			e.Types = append(e.Types, 1)
		} else if endpointType == "Transaction" {
			e.Types = append(e.Types, 2)
		} else if endpointType == "Block" {
			e.Types = append(e.Types, 3)
		} else if endpointType == "RequestMissingTransaction"{
			e.Types = append(e.Types,4)
		} else {
			fmt.Println("Bad endpoint type!")
			os.Exit(12)
		}
	}
}

func (e *Endpoint)GetEndpointTypes() []string {
	res := make([] string, len(e.Types))
	for _, endpointType := range e.Types {
		if endpointType == 1 {
			res = append(res, "FailureDetection")
		} else if endpointType == 2 {
			res = append(res, "Transaction")
		} else if endpointType == 3 {
			res = append(res, "Block")
		} else if endpointType == 4 {
			res = append(res, "RequestMissingTransaction")
		}
	}
	return  res
}
