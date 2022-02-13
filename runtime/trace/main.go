package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	const number_of_random_bytes = 1024 * 1024 * 1024
	data := make([][]byte, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		data[i] = make([]byte, number_of_random_bytes)
		n, err := rand.Read(data[i])
		if err != nil {
			log.Fatal(err)
		}
		if n != number_of_random_bytes {
			log.Fatal("Number of random bytes read is not 4096 for thread", fmt.Sprint(i))
		}
	}

	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rand.Seed(time.Now().UnixNano())

	log.Println("Start trace now ...")
	trace.Start(f)
	defer func() {
		trace.Stop()
		log.Println("Tracing stopped.")
	}()

	var dummy byte = 0xFF
	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(data []byte) {
			defer wg.Done()
			hash := sha256.Sum256(data)
			dummy &= hash[0]
		}(data[i])
		//for i := 0; i < len(hash); i++ {
		//	fmt.Fprintf(f, "%02x", hash[i])
		//}
		//fmt.Fprintf(f, "\n")
	}
	wg.Wait()
}
