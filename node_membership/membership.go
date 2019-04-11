package node_membership

import (
	"fmt"
	"sort"
	"sync"
)

type Membership struct {
	List           [] Entry
	ListMutex      sync.Mutex
}


func (m *Membership) PrintContent() {
	fmt.Println("--------------------")
	fmt.Println("Membership content:")
	for _, v := range m.List {
		fmt.Println("Name:", v.Name, ", ip:", v.IpAddress,", type:",v.DecodeEntryType(), ", incarn:",v.Incarnation, ", lastTime:", v.LastUpdatedTime)
	}
	fmt.Println("=====================")
}
/*
	@param ipAddress string
	@param incarnation int
	@param entryType int
	Invoke when the server receives response from ping.  Update the membershipList
 */
func (m *Membership) UpdateNode(entry Entry) {
	m.ListMutex.Lock()
	for i, elem := range m.List {
		if elem.IpAddress == entry.IpAddress  {
			if entry.EntryType == 0 {
				//new entry is Alive
				if entry.Incarnation > elem.Incarnation {
					if m.List[i].EntryType == 0  {
						m.List[i].EntryType = 0
						m.List[i].Incarnation = entry.Incarnation
					} else if m.List[i].EntryType == 1 {
						m.List[i].EntryType = 0
						m.List[i].Incarnation = entry.Incarnation
					}

				}
			}
			if entry.EntryType == 1 {
				//new entry is Suspect
				if entry.Incarnation > elem.Incarnation && elem.EntryType == 1 {
					m.List[i].EntryType = 1
					m.List[i].Incarnation = entry.Incarnation
				}

				if entry.Incarnation >= elem.Incarnation && elem.EntryType == 0 {
					m.List[i].EntryType = 1
					m.List[i].Incarnation = entry.Incarnation
				}
			}

			if entry.EntryType == 2 {
				//new entry is Failure
				m.List[i].EntryType = 2
				m.List[i].Incarnation = entry.Incarnation
			}
			m.ListMutex.Unlock()
			return
		}
	}
	m.ListMutex.Unlock()
	m.AddNewNode(entry)
	m.SortMembership()
}
func (m *Membership) SortMembership(){
	m.ListMutex.Lock()
	defer m.ListMutex.Unlock()
	sort.Slice(m.List, func(i, j int) bool {
		key1 := fmt.Sprintf("%s%d",m.List[i].IpAddress,m.List[i].InitialTimeStamp )
		key2 := fmt.Sprintf("%s%d",m.List[j].IpAddress,m.List[j].InitialTimeStamp )
		return  key1<key2
	})
}
func (m *Membership) UpdateNode2(ipAddress string, entryType int, lastUpdatedTime int64) {
	m.ListMutex.Lock()
	defer m.ListMutex.Unlock()

	for i, elem := range m.List {
		if elem.IpAddress == ipAddress {
			m.List[i].EntryType = entryType
			m.List[i].LastUpdatedTime = lastUpdatedTime
			return
		}
	}
}

/*
	@param ipAddress string
	@param initialTimeStamp int64
	Invoke when the server receives response from ping.  Update the membershipList
 */
func (m *Membership) AddNewNode(entry Entry) {
	if m.ContainsNode(entry) {
		panic("ip address is already in the list")
	}

	m.ListMutex.Lock()
	m.List = append(m.List, entry)
	m.ListMutex.Unlock()
}

/*
	@param ipAddress string
	@param initialTimeStamp int64
	Invoke when the server receives response from ping.  Update the membershipList and the disseminateList
 */
func (m *Membership) RemoveNode(ipAddress string, initialTimeStamp int64) {
	m.ListMutex.Lock()
	defer m.ListMutex.Unlock()
	for ind, elem := range m.List {
		if elem.IpAddress == ipAddress && elem.InitialTimeStamp == initialTimeStamp {
			m.List = append(m.List[:ind], m.List[ind+1:]...)
			return
		}
	}
}

func (m *Membership) ContainsNode(entry Entry) bool {
	m.ListMutex.Lock()
	defer m.ListMutex.Unlock()

	for _, elem := range m.List {
		if elem.IpAddress == entry.IpAddress && elem.Name == entry.Name {
			return true
		}
	}
	return false
}