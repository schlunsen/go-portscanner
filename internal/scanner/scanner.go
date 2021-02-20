package scanner

import (
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/cheggaaa/pb"
)

type Scanner struct {
	Hostname  string
	Startport int16
	Endport   int16
}

func NewScanner(hostname string, startport int16, endport int16) *Scanner {
	s := new(Scanner)
	s.Hostname = hostname
	s.Startport = startport
	s.Endport = endport
	return s
}

func worker(hostname string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", hostname, p)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func (s Scanner) Scan() {
	ports := make(chan int, 100)
	// this channel will receive results of scanning
	results := make(chan int)
	var openports []int

	var bar *pb.ProgressBar = pb.StartNew(int(s.Endport) - int(s.Startport))

	// create a pool of workers
	for i := 0; i < cap(ports); i++ {
		go worker(s.Hostname, ports, results)
	}

	// send ports to be scanned
	go func() {
		for i := int(s.Startport); i <= int(s.Endport); i++ {
			ports <- i
		}
	}()

	for i := s.Startport; i <= s.Endport; i++ {
		port := <-results
		bar.Increment()
		if port != 0 {
			fmt.Printf("%d open\n", port)
			openports = append(openports, port)
		}
	}

	// After all the work has been completed, close the channels
	close(ports)
	close(results)
	// sort open port numbers
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}

}
