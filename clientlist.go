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

	cl.clientList = append(cl.clientList, o)
	cl.total = len(cl.clientList)
}

func (cl *ClientList) findFree() (IClient, bool) {

	cl.mu.Lock()
	defer cl.mu.Unlock()

	for i := range cl.clientList {

		o := cl.clientList[(i+cl.lastI)%cl.total]

		if !o.isLock() {
			cl.lastI = i + 1
			return o.lock(), true
		}
	}

	return nil, false
}
