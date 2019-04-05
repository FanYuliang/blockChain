package endpoints

import "mp2/node_membership"

type FailureDetectionMeta struct {
	Type             int // 0: join, 1: ping, 2: ack 3: leave
	Record           []node_membership.Entry
	InitialTimeStamp int64
	IpAddress        string
}

func (f *FailureDetectionMeta) EncodeFailureDetectionActionType(endpointType string) int {
	if endpointType == "Join" {
		return 1
	} else if endpointType == "Ping" {
		return 2
	} else if endpointType == "Ack" {
		return 3
	} else if endpointType == "Quit" {
		return 4
	}
	return -1
}