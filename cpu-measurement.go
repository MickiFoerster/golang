package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type cpustat struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	total   uint64
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
			var stat cpustat
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
				default:
				}
				totalticks += val
			}
			stat.total = totalticks
			return &stat, nil
		}
	}
	return nil, fmt.Errorf("Could not fine line with prefix cpu")
}

func main() {
	stat, err := getCPUSample()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stat)
	/*
		idle0, total0 := getCPUSample()
		time.Sleep(3 * time.Second)
		idle1, total1 := getCPUSample()

		idleTicks := float64(idle1 - idle0)
		totalTicks := float64(total1 - total0)
		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

		fmt.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
	*/
}
