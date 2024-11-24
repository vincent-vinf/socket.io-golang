package socketio

import (
	"sync"
)

type Namespace struct {
	name         string
	sockets      connections
	onConnection connectionEvent
}

func newNamespace(name string) *Namespace {
	return &Namespace{
		name: name,
		sockets: connections{
			conn: make(map[string]*Socket),
		},
		onConnection: connectionEvent{
			list: make(map[string][]connectionEventCallback),
		},
	}
}

func (nps *Namespace) OnConnection(fn connectionEventCallback) {
	nps.onConnection.set("connection", fn)
}

func (nps *Namespace) Emit(event string, agrs ...interface{}) error {
	for _, socket := range nps.sockets.all() {
		socket.Emit(event, agrs...)
	}
	return nil
}

type namespaces struct {
	sync.RWMutex
	list map[string]*Namespace
}

func (n *namespaces) create(name string) *Namespace {
	n.Lock()
	n.list[name] = newNamespace(name)
	ret := n.list[name]
	n.Unlock()
	return ret
}

func (n *namespaces) all() []*Namespace {
	n.RLock()
	ret := make([]*Namespace, 0)
	for _, nps := range n.list {
		ret = append(ret, nps)
	}
	n.RUnlock()
	return ret
}

func (n *namespaces) get(name string) *Namespace {
	n.RLock()
	defer n.RUnlock()
	ret, ok := n.list[name]
	if !ok {
		return nil
	}
	return ret
}
