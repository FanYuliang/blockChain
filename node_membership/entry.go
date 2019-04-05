package node_membership

import "log"

type Entry struct {
	Name             string
	IpAddress        string
	InitialTimeStamp int64
	Incarnation      int
	EntryType        int //0 for alive, 1 for suspected, 2 for failed
	LastUpdatedTime  int64
}

func (e *Entry) EncodeEntryType(entryType string) int {
	if entryType == "alive" {
		return 0
	} else if entryType == "suspected" {
		return 1
	} else if entryType == "failed" {
		return 2
	}
	log.Fatalln("Fail to encode entry type ", entryType)
	return -1
}
