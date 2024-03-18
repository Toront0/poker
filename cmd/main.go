package main

import (
	"github.com/Toront0/poker/internal/api"

)

func main() {
	

	server := api.NewServer(":3000")

	server.Run()

}

// package main

// import (
// 	"sync"
// 	"sync/atomic"
// 	"time"
// 	"fmt"

// )

// func main() {
// 	AddMutex()
// 	// AddAtomic()
// }

// func AddMutex() {
// 	start := time.Now()

// 	var (
// 		counter int64
// 		wg sync.WaitGroup
// 		mu sync.Mutex
// 	)

// 	wg.Add(1000)

// 	for i := 0; i < 1000; i++ {

// 		go func() {

// 			defer wg.Done()

// 			mu.Lock()

// 			counter++

// 			mu.Unlock()

// 		}()

// 	}
// 	wg.Wait()

// 	fmt.Println("counter", counter)
// 	fmt.Println("with mutex: ", time.Now().Sub(start).Seconds())

// }

// func AddAtomic() {
// 	start := time.Now()

// 	var (
// 		counter int64
// 		wg sync.WaitGroup
// 	)

// 	wg.Add(1000)

// 	for i := 0; i < 1000; i++ {

// 		go func() {

// 			defer wg.Done()

// 			atomic.AddInt64(&counter, 1)

// 		}()

// 	}
// 	wg.Wait()

// 	fmt.Println("counter", counter)
// 	fmt.Println("with atomic: ", time.Now().Sub(start).Seconds())
// }