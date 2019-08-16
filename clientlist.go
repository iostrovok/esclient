package esclient

import (
	"sync"
)

type ClientList struct {
	mu sync.RWMutex

	clientList []*oneClient
	lastI      int
	total      int
}

func NewClientList() *ClientList {
	return &ClientList{
		clientList: []*oneClient{},
	}
}

func (cl *ClientList) appendTo(o *oneClient) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	if o.ConnError() == nil {
		cl.clientList = append(cl.clientList, o)
		cl.total = len(cl.clientList)
	}
}

func (cl *ClientList) findFree() (IClient, bool) {

	cl.mu.Lock()
	defer cl.mu.Unlock()

	for i := range cl.clientList {

		j := (i + cl.lastI) % cl.total

		o := cl.clientList[j]

		if o.ConnError() != nil {
			if len(cl.clientList) > j {
				cl.clientList = append(cl.clientList[:j], cl.clientList[j+1:]...)
				cl.total = len(cl.clientList)
				cl.lastI--
			}
			continue
		}

		if !o.isLock() {
			cl.lastI = i + 1
			return o.lock(), true
		}
	}

	return nil, false
}
