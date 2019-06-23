package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getMemSample() float64 {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatalln("error while opening file /proc/meminfo:", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var totalKB float64
	var availKB float64
loop:
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			total := false
			free := false
			avail := false
			switch {
			case strings.HasPrefix(line, "MemTotal:"):
				line = line[len("MemTotal: ") : len(line)-3]
				total = true
			case strings.HasPrefix(line, "MemFree:"):
				line = line[len("MemFree: ") : len(line)-3]
				free = true
			case strings.HasPrefix(line, "MemAvailable:"):
				line = line[len("MemAvailable: ") : len(line)-3]
				avail = true
			default:
				break loop
			}

			line = strings.Trim(line, " ")
			val, err := strconv.ParseUint(line, 10, 64)
			if err != nil {
				fmt.Println("Parse error: ", line, err)
				continue
			}

			switch {
			case total:
				totalKB = float64(val)
			case free:
				_ = val
			case avail:
				availKB = float64(val)
			default:
				log.Fatal("unexpected error")
			}
		}
	}
	memUsed := totalKB - availKB
	return memUsed * 100 / totalKB
}

func getMemUsage() <-chan float64 {
	c := make(chan float64)
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				c <- getMemSample()
			}
		}
	}()

	return c
}

func main() {
	c := getMemUsage()
	for {
		select {
		case mu := <-c:
			fmt.Println(mu)
		}
	}
}
