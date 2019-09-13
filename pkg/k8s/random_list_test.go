package k8s

import (
	"fmt"
	"testing"
)

func TestRandomList_Random(t *testing.T) {
	rl := NewRandomList()
	for i := 0; i < 10; i++ {
		rl.Add(fmt.Sprintf("%d", i))
	}

	m := make(map[string]int)
	for i := 0; i < 100; i++ {
		rnd := rl.Random()
		_, ok := m[rnd]
		if ok {
			m[rnd]++
		} else {
			m[rnd] = 1
		}
	}

	for key, value := range m {
		fmt.Println(key, value)
	}
	fmt.Println("remove")
	for i := 0; i < 5; i++ {
		rl.Remove(fmt.Sprintf("%d", i))
	}

	for i := 0; i < 100; i++ {
		rnd := rl.Random()
		_, ok := m[rnd]
		if ok {
			m[rnd]++
		} else {
			m[rnd] = 1
		}
	}

	for key, value := range m {
		fmt.Println(key, value)
	}
}
