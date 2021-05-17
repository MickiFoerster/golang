// Concurrent Monte-Carlo simulation.
// Bottleneck seems to be the random number generator
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	const (
		e           = math.E
		experiments = 100_000_000
	)

	fmt.Printf("Estimating e with %d experiment(s).\n\n", experiments)

	N := runtime.GOMAXPROCS(0)

	var overall_tries int = 0
	var mtx sync.Mutex
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		seed := time.Now().UTC().UnixNano() % 100_000
		go func(s int64) {
			defer wg.Done()
			rand.Seed(s)
			fmt.Fprintf(os.Stderr, "This Go routine starts Monte-Carlo simulation with seed: %d\n", s)
			var local_tries int
			for j := 0; j < experiments/N; j++ {
				var (
					sum         float64
					num2Success int
				)

				for sum <= 1 {
					n := rand.Float64()
					sum += n
					num2Success++
				}
				local_tries += num2Success
			}
			mtx.Lock()
			overall_tries += local_tries
			mtx.Unlock()
		}(seed)
	}
	wg.Wait()

	expected := float64(overall_tries) / float64(experiments)
	errorPct := 100.0 * math.Abs(expected-e) / e

	fmt.Printf("Expected vale: %9f \n", expected)
	fmt.Printf("e: %9f \n", e)
	fmt.Printf("Error: %9f%%\n", errorPct)
}
