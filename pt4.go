package main

import (
	"fmt"
	"sync"
)

func new(i int) int {
	return i 
}

func main() {
	var wg sync.WaitGroup
	var m sync.Mutex
	for i := 10; i >= 0; i--{
		wg.Add(1)
		go func(i int ){
			m.Lock()
			defer m.Unlock()
			defer wg.Done()
			fmt.Println(i)
		} (i)
	}
	wg.Wait()
}