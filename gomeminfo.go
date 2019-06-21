package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatalln("error while opening file /proc/meminfo:", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			match := strings.HasPrefix(line, "MemTotal:") ||
				strings.HasPrefix(line, "MemFree:") ||
				strings.HasPrefix(line, "MemAvailable:")
			if match {
				fmt.Println(line)
			}
		}
	}
}
