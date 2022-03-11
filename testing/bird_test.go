package testing

import (
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	c := make(chan int)
	go func() {
		fmt.Println("ready to send in g1")
		c <- 1
		fmt.Println("send 1 to chan")
		fmt.Println("goroutine start sleep 1 second")
		time.Sleep(time.Second)
		fmt.Println("goroutine end sleep")
		c <- 2
		fmt.Println("send 2 to chan")
	}()

	fmt.Println("main thread start sleep 1 second")
	time.Sleep(time.Second)
	fmt.Println("main thread and sleep")
	i := <-c
	fmt.Printf("receive %d\n", i)
	i = <-c
	fmt.Printf("receive %d\n", i)
	time.Sleep(time.Second)
}
