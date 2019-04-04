package server

import "mp2/endpoints"

func (s *Server) getFailureDetectionEndpointMetadata(ActionType string) endpoints.FailureDetectionMeta {
	num := int(float32(len(s.MembershipList.List))*0.3)

	if num < 1 {
		num = 1
	}
	listToSend := s.getMemebershipSubset(num)

	fEndpoint := endpoints.FailureDetectionMeta{
		endpoints.EncodeFailureDetectionActionType(ActionType),
		listToSend,
		s.InitialTimeStamp,
		s.MyAddress}
	return fEndpoint
}

func (s *Server) getTransactionEndpointMetadata() endpoints.TransactionMeta {
	transactionToSend := s.getTransactSubset()
	tEndpoint := endpoints.TransactionMeta{transactionToSend}
	return tEndpoint
}
