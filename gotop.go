// Measure CPU usage
// Hints can be found in https://www.idnt.net/en-US/kb/941772

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

type cpustat struct {
	user       uint64
	nice       uint64
	system     uint64
	idle       uint64
	iowait     uint64
	irq        uint64
	softirq    uint64
	steal      uint64
	guest      uint64
	guest_nice uint64
	total      uint64
}

func (stat *cpustat) String() string {
	s := fmt.Sprintf("%10s %10v\n", "user:", stat.user)
	s += fmt.Sprintf("%10s %10v\n", "nice:", stat.nice)
	s += fmt.Sprintf("%10s %10v\n", "system:", stat.system)
	s += fmt.Sprintf("%10s %10v\n", "idle:", stat.idle)
	s += fmt.Sprintf("%10s %10v\n", "iowait:", stat.iowait)
	s += fmt.Sprintf("%10s %10v\n", "irq:", stat.irq)
	s += fmt.Sprintf("%10s %10v\n", "softirq:", stat.softirq)
	s += fmt.Sprintf("%10s %10v\n", "total:", stat.total)
	return s
}

func getCPUSample() (*cpustat, error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			stat := new(cpustat)
			var totalticks uint64
			totalticks = 0
			for i := 1; i < len(fields); i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
					continue
				}
				switch i {
				case 1:
					stat.user = val
				case 2:
					stat.nice = val
				case 3:
					stat.system = val
				case 4:
					stat.idle = val
				case 5:
					stat.iowait = val
				case 6:
					stat.irq = val
				case 7:
					stat.softirq = val
				case 8:
					stat.steal = val
				case 9:
					stat.guest = val
				case 10:
					stat.guest_nice = val
				default:
				}
				totalticks += val
			}
			stat.total = totalticks
			return stat, nil
		}
	}
	return nil, fmt.Errorf("Could not fine line with prefix cpu")
}

func getCpuUsage(stat2, stat1 *cpustat) float64 {
	idleTicks := float64(stat2.idle - stat1.idle)
	waiting := idleTicks
	totalDelta := float64(stat2.total - stat1.total)
	cpuUsed := totalDelta - waiting
	return 100 * cpuUsed / totalDelta
}

func loop(ch chan<- float64) {
	stat1, err := getCPUSample()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)
	stat2, err := getCPUSample()
	if err != nil {
		log.Fatal(err)
	}

	cpuUsage := getCpuUsage(stat2, stat1)

	stat1 = stat2
	ch <- cpuUsage

	for {
		select {
		case <-time.After(500 * time.Millisecond):
			stat2, err := getCPUSample()
			if err != nil {
				log.Fatal(err)
			}

			cpuUsage = getCpuUsage(stat2, stat1)

			stat1 = stat2
			ch <- cpuUsage
		}
	}
}

func main() {
	ch := make(chan float64)
	go loop(ch)

	for {
		select {
		case cpuUsage := <-ch:
			fmt.Printf("%s: CPU usage is %f%%\r", time.Now().String()[:19], cpuUsage)
		}
	}
}
