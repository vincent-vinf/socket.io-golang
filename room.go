package socketio

import (
	"sync"
)

type Room struct {
	name    string
	sockets connections
}

func newRoom(name string) *Room {
	return &Room{
		name: name,
		sockets: connections{
			conn: make(map[string]*Socket),
		},
	}
}

func (room *Room) Emit(event string, agrs ...interface{}) error {
	for _, socket := range room.sockets.all() {
		socket.Emit(event, agrs...)
	}
	return nil
}

type roomNames struct {
	sync.RWMutex
	list []string
}

func (l *roomNames) set(name string) {
	l.Lock()
	l.list = append(l.list, name)
	l.Unlock()
}

func (l *roomNames) delete(name string) int {
	l.Lock()
	defer l.Unlock()
	for index, n := range l.list {
		if n == name {
			l.list = append(l.list[:index], l.list[index+1:]...)
			return index
		}
	}
	return -1
}

func (l *roomNames) all() []string {
	l.RLock()
	defer l.RUnlock()
	return append([]string{}, l.list...)
}

// func (l *roomNames) get(event string) []eventCallback {
// 	l.RLock()
// 	defer l.RUnlock()
// 	if _, ok := l.list[event]; !ok {
// 		return make([]eventCallback, 0)
// 	}
// 	ret := make([]eventCallback, 0)
// 	ret = append(ret, l.list[event]...)
// 	return ret
// }

type rooms struct {
	sync.RWMutex
	list map[string]*Room
}

func (n *rooms) create(name string) *Room {
	n.Lock()
	ret, ok := n.list[name]
	if !ok {
		ret = newRoom(name)
		n.list[name] = ret
	}
	n.Unlock()
	return ret
}

// func (n *rooms) get(name string) *Room {
// 	n.RLock()
// 	defer n.RUnlock()
// 	ret, ok := n.list[name]
// 	if !ok {
// 		return nil
// 	}
// 	return ret
// }
