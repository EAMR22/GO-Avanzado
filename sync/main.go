package main

import (
	"fmt"
	"sync"
)

var (
	balance int = 100
)

func Deposit(amount int, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.Lock()
	b := balance
	balance = b + amount
	lock.Unlock()
}

func Balance(lock *sync.RWMutex) int {
	b := balance
	return b
}

func main() {
	var wg sync.WaitGroup
	var lock sync.RWMutex // Es un look de escritura y lectura.
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go Deposit(i*100, &wg, &lock)
	}
	wg.Wait()
	fmt.Println(Balance(&lock))
}