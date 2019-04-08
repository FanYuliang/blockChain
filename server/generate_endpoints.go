package server

import (
	"mp2/blockchain"
	"mp2/endpoints"
)

func (s *Server) getFailureDetectionEndpointMetadata(ActionType string) endpoints.FailureDetectionMeta {
	num := int(float32(len(s.MembershipList.List)) * 0.3)

	if num < 1 {
		num = 1
	}
	listToSend := s.getMemebershipSubset(num)

	var fEndpoint endpoints.FailureDetectionMeta
	fEndpoint.Type = fEndpoint.EncodeFailureDetectionActionType(ActionType)
	fEndpoint.Record = listToSend
	fEndpoint.InitialTimeStamp = s.InitialTimeStamp
	fEndpoint.IpAddress = s.MyAddress
	return fEndpoint
}

func (s *Server) getTransactionEndpointMetadata() endpoints.TransactionMeta {
	transactionToSend := s.Transactions.GetTransactionToCommit(20)
	tEndpoint := endpoints.TransactionMeta{Tx: transactionToSend}
	//fmt.Println("tEndpoint: ", tEndpoint)
	return tEndpoint
}

func (s *Server) getRequestMissingBlockMeta(missingBlockID string) endpoints.RequestMissingBlockMeta {
	missingBlockMeta := endpoints.RequestMissingBlockMeta{MissingBlockID: missingBlockID, RequesterIPaddr: s.MyAddress}
	return missingBlockMeta
}

func (s *Server) getBlockMeta(b blockchain.Block) endpoints.BlockMeta {
	blockmeta := endpoints.BlockMeta{Block: b}
	return blockmeta
}
