package main

import (
	"crypto/sha256"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	data := make([]byte, 4096)
	n, err := rand.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	if n != 4096 {
		log.Fatal("Number of random bytes read is not 4096")
	}

	f, err := os.Create("cpuprofile.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rand.Seed(time.Now().UnixNano())

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	var dummy byte = 0xFF
	for i := 0; i < 1_000_000; i++ {
		hash := sha256.Sum256(data)
		dummy &= hash[0]
		//for i := 0; i < len(hash); i++ {
		//	fmt.Fprintf(f, "%02x", hash[i])
		//}
		//fmt.Fprintf(f, "\n")
	}
}
