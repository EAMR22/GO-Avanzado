package main

import (
	"fmt"
	"sync"
)

var (
	balance int = 100
)

func Deposit(amount int, wg *sync.WaitGroup, lock *sync.Mutex) {
	defer wg.Done() // Decrementa el contador en 1.
	lock.Lock()     // Bloquea el programa, hasta que termine las modificaciones.
	b := balance
	balance = b + amount
	lock.Unlock() // Desbloquea el programa, al haber terminado de hacer las modificaciones.
}

func Balance() int {
	b := balance
	return b
}

func main() {
	var wg sync.WaitGroup
	var lock sync.Mutex
	for i := 1; i <= 5; i++ {
		wg.Add(1) // Incrementa el contador en 1.
		go Deposit(i*100, &wg, &lock)
	}
	wg.Wait() // Bloquea el programa.
	fmt.Println(Balance())
}
