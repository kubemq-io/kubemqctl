package k8s

import "sync"
import "math/rand"
import "time"

type RandomList struct {
	list []string
	sync.Mutex
}

func NewRandomList() *RandomList {
	rand.Seed(time.Now().UnixNano())
	r := &RandomList{
		list:  nil,
		Mutex: sync.Mutex{},
	}
	return r
}

func (rl *RandomList) Add(item string) {
	rl.Lock()
	defer rl.Unlock()
	rl.list = append(rl.list, item)
}
func (rl *RandomList) Remove(item string) {
	rl.Lock()
	defer rl.Unlock()
	removeIndex := -1
	for index, value := range rl.list {
		if value == item {
			removeIndex = index
			break
		}
	}
	rl.list[removeIndex] = rl.list[len(rl.list)-1] // Copy last element to index i.
	rl.list[len(rl.list)-1] = ""                   // Erase last element (write zero value).
	rl.list = rl.list[:len(rl.list)-1]             // Truncate slice.
}

func (rl *RandomList) Random() string {
	rl.Lock()
	defer rl.Unlock()
	if len(rl.list) == 0 {
		return ""
	}
	return rl.list[rand.Intn(len(rl.list))]
}
func (rl *RandomList) Len() int {
	rl.Lock()
	defer rl.Unlock()
	return len(rl.list)
}
