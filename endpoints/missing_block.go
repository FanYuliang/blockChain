package endpoints

import "mp2/blockchain"

type RequestMissingBlockMeta struct{
	MissingBlockID    	string
	RequesterIPaddr 	string
}

type SendMissingBlockMeta struct{
	MissingBlock    blockchain.Block
	RequesterIPaddr string
}