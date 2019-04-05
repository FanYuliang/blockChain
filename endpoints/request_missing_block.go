package endpoints

import "mp2/blockchain"

type RequestMissingBlocknMeta struct{
	Type 			int // 0 receive, 1 reply
	RequesterIPaddr string
	MissingBlock 	blockchain.Block
}
