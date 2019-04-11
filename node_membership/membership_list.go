package node_membership

import (
	"fmt"
	"mp2/utils"
	"sort"
	"sync"
	"time"
)
// MembershipList the set of Items
type MembershipList struct {
	MyIPAddr string
	items []Entry
	lock  sync.RWMutex
}

func (m *MembershipList) Copy() *MembershipList {
	m.lock.RLock()
	defer m.lock.RUnlock()
	res := new(MembershipList)
	for _, v := range m.items {
		res.items = append(res.items, v)
	}
	return res
}

func (m *MembershipList) Size() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.items)
}

func (m *MembershipList) PrintContent() {
	m.lock.RLock()
	defer m.lock.RUnlock()
	fmt.Println("--------------------")
	fmt.Println("Membership content:")
	for _, v := range m.items {
		fmt.Println("Name:", v.Name, ", ip:", v.IpAddress,", type:",v.DecodeEntryType(), ", incarn:",v.Incarnation, ", lastTime:", v.LastUpdatedTime)
	}
	fmt.Println("=====================")
}

func (m *MembershipList) UpdateNode(entry Entry) {
	m.lock.Lock()
	for i, elem := range m.items {
		if elem.IpAddress == entry.IpAddress  {
			if entry.EntryType == 0 {
				//new entry is Alive
				if entry.Incarnation > elem.Incarnation {
					if m.items[i].EntryType == 0  {
						m.items[i].EntryType = 0
						m.items[i].Incarnation = entry.Incarnation
					} else if m.items[i].EntryType == 1 {
						m.items[i].EntryType = 0
						m.items[i].Incarnation = entry.Incarnation
					}
				}
			}
			if entry.EntryType == 1 {
				//new entry is Suspect
				if entry.Incarnation > elem.Incarnation && elem.EntryType == 1 {
					m.items[i].EntryType = 1
					m.items[i].Incarnation = entry.Incarnation
				} else if entry.Incarnation >= elem.Incarnation && elem.EntryType == 0 {
					m.items[i].EntryType = 1
					m.items[i].Incarnation = entry.Incarnation
				} else if entry.IpAddress == m.MyIPAddr {
					//prove myself!
					m.items[i].EntryType = 0
					m.items[i].Incarnation += 1
				}
			}

			if entry.EntryType == 2 {
				//new entry is Failure
				m.items[i].EntryType = 2
				m.items[i].Incarnation = entry.Incarnation

				if entry.IpAddress == m.MyIPAddr {
					fmt.Println("我tmd还活着啊！！！ from ", m.items[i].Name)
				}
			}
			m.lock.Unlock()
			return
		}
	}
	m.lock.Unlock()
	m.AddNewNode(entry)

	//sort membership
	m.lock.Lock()
	sort.Slice(m.items, func(i, j int) bool {
		key1 := fmt.Sprintf("%s%d", m.items[i].IpAddress, m.items[i].InitialTimeStamp )
		key2 := fmt.Sprintf("%s%d", m.items[j].IpAddress, m.items[j].InitialTimeStamp )
		return  key1<key2
	})
	m.lock.Unlock()
}

func (m *MembershipList) UpdateNode2(ipAddress string, entryType int, lastUpdatedTime int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for i, elem := range m.items {
		if elem.IpAddress == ipAddress {
			m.items[i].EntryType = entryType
			m.items[i].LastUpdatedTime = lastUpdatedTime
			return
		}
	}
}

func (m *MembershipList) SetUpdatedTime(ipAddress string,lastUpdatedTime int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for i, elem := range m.items {
		if elem.IpAddress == ipAddress {
			m.items[i].LastUpdatedTime = lastUpdatedTime
			return
		}
	}
}

/*
	@param ipAddress string
	@param initialTimeStamp int64
	Invoke when the server receives response from ping.  Update the membershipList
 */
func (m *MembershipList) AddNewNode(entry Entry) {

	if m.ContainsNode(entry) {
		panic("ip address is already in the list")
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.items = append(m.items, entry)
}

/*
	@param ipAddress string
	@param initialTimeStamp int64
	Invoke when the server receives response from ping.  Update the membershipList and the disseminateList
 */
func (m *MembershipList) RemoveNode(ipAddress string, initialTimeStamp int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for ind, elem := range m.items {
		if elem.IpAddress == ipAddress && elem.InitialTimeStamp == initialTimeStamp {
			m.items = append(m.items[:ind], m.items[ind+1:]...)
			return
		}
	}
}

func (m *MembershipList) ContainsNode(entry Entry) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, elem := range m.items {
		if elem.IpAddress == entry.IpAddress && elem.Name == entry.Name {
			return true
		}
	}
	return false
}


func (m *MembershipList) GetEntryByIndex(index int) Entry {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.items[index]
}

func (m *MembershipList) GetAll() []Entry {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.items
}

func (m *MembershipList) GetNonFailureMembershipSize() {
	m.lock.RLock()
	defer m.lock.RUnlock()

	size := 0
	for _, v := range m.items {
		if v.EntryType != 2 {
			size += 1
		}
	}
	//fmt.Println("Non failure membership size: ", size)
}

func (m *MembershipList) FindEntryByIP(targetIP string) int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for ind, entry := range m.items {
		if targetIP == entry.IpAddress {
			return ind
		}
	}

	fmt.Println("Fail to find self in membership list.")
	return -1
}

func (m *MembershipList) GetMemebershipSubset(subsetNum int) []Entry {
	m.lock.RLock()
	defer m.lock.RUnlock()
	tempArr := utils.Arange(0, len(m.items), 1)
	shuffledArr := utils.Shuffle(tempArr)
	var res [] Entry

	for _, v := range m.items {
		if v.IpAddress == m.MyIPAddr {
			res = append(res, v)
			break
		}
	}
	for _, v := range shuffledArr {
		if len(res) > subsetNum {
			break
		}
		if m.items[v].IpAddress != m.MyIPAddr {
			res = append(res, m.items[v])
		}
	}
	return res
}

func (m *MembershipList) CheckMembershipList(detectionTimeout int64, suspectTimeout int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	currTime := time.Now().Unix()
	//check if any process is MembershipList or failed
	for i := len(m.items) - 1; i >= 0; i-- {
		entry := m.items[i]
		if entry.EntryType == 0 && currTime-entry.LastUpdatedTime >= detectionTimeout && entry.LastUpdatedTime != 0 {
			//alive now but passed detection timeout
			m.items[i].LastUpdatedTime = 0
			m.items[i].EntryType = 1
		} else if entry.EntryType == 1 && currTime-entry.LastUpdatedTime >= suspectTimeout && entry.LastUpdatedTime != 0 {
			//suspected now but passed suspected timeout
			m.items[i].EntryType = 2
		}
	}
}

func (m *MembershipList) GetEntryByIpAddress(targetIP string) Entry {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, entry := range m.items {
		if targetIP == entry.IpAddress {
			return entry
		}
	}

	fmt.Println("Fail to find self in membership list.")
	return Entry{}
}

func (m *MembershipList) UpdateLastUpdatedTimeWithIP(targetIP string, lastUpdatedTime int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for i, entry := range m.items {
		if targetIP == entry.IpAddress {
			m.items[i].LastUpdatedTime = lastUpdatedTime
			break
		}
	}
}