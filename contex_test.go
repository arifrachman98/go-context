package gocontex

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contA := context.Background()

	contB := context.WithValue(contA, "b", "B")
	contC := context.WithValue(contA, "c", "C")

	contD := context.WithValue(contB, "d", "D")
	contE := context.WithValue(contB, "e", "E")

	contF := context.WithValue(contC, "f", "F")
	contG := context.WithValue(contF, "g", "G")

	fmt.Println(contA)
	fmt.Println(contB)
	fmt.Println(contC)
	fmt.Println(contD)
	fmt.Println(contE)
	fmt.Println(contF)
	fmt.Println(contG)

	fmt.Println(contF.Value("f")) //data dapat muncul karna key f pada context F memiliki value F
	fmt.Println(contF.Value("c")) //data dapat muncul karna key c pada context F memiliki value F dari parent key c
	fmt.Println(contF.Value("b")) //data tidak dapat muncul karna key b pada context F tidak memiliki value B dan tidak memiliki parent yang memiliki value B
}

func CreateCounterCancel(ctx context.Context) chan int {
	destination := make(chan int)
	go func() {
		defer close(destination)
		count := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- count
				count++
			}
		}
	}()
	return destination
}

func CreateCounterTimeout(ctx context.Context) chan int {
	destination := make(chan int)
	go func() {
		defer close(destination)
		count := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- count
				count++
				time.Sleep(1 * time.Second) //simulasi slow response
			}
		}
	}()
	return destination
}

func CreateCounterDeadline(ctx context.Context) chan int {
	destination := make(chan int)
	go func() {
		defer close(destination)
		count := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- count
				count++
				time.Sleep(1 * time.Second) //simulasi slow response
			}
		}
	}()
	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background() //init background process
	ctx, cancel := context.WithCancel(parent)
	dest := CreateCounterCancel(ctx)

	for n := range dest {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}
	cancel()

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background() //init background process
	ctx, cancel := context.WithTimeout(parent, 3*time.Second)
	dest := CreateCounterTimeout(ctx)
	defer cancel()

	//forever looping, will stop with timeout
	for n := range dest {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithTimeoutDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background() //init background process
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(4*time.Second))
	dest := CreateCounterDeadline(ctx)
	defer cancel()

	//forever looping, will stop with deadline time
	for n := range dest {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
