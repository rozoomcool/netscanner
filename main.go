package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// hosts := make(chan string)
	ports := make(chan int)
	defer close(ports)
	// defer close(hosts)

	// getEnabledHosts(hosts)

	// for host := range hosts {
	checkAddressPorts("localhost", ports)
	// }

	// for host := range hosts {
	// 	fmt.Printf("HOST::%v\n", host)
	for i := range ports {
		if i != 0 {
			fmt.Printf("\tport %d is open\n", i)
		}
	}

}

func getEnabledHosts(hosts chan string) {
	baseIP := getPrefId()

	go func() {
		for i := 1; i <= 255; i++ {
			cmd := exec.Command("ping", "-c", "1", "-W", "1", fmt.Sprintf("%v:%d", baseIP, i))
			if error := cmd.Run(); error != nil {
				continue
			}

			if error := cmd.Wait(); error != nil {
				continue
			}

			hosts <- fmt.Sprintf("%v:%d", baseIP, i)
		}
	}()
}

func getPrefId() string {
	addr, err := GetOutboundIP()
	if err != nil {
		fmt.Println(err)
	}
	ip := addr.String()

	ipArr := strings.Split(ip, ".")
	baseIP := strings.Join(ipArr[0:len(ipArr)-2], ".")

	return baseIP
}

func checkAddressPorts(port string, result chan int) {
	go func() {
		for i := 1; i <= 40000; i++ {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%d", port, i), time.Second*3)
			if err != nil {
				result <- 0
			} else {
				conn.LocalAddr()
				result <- i
			}
		}
	}()
}

func GetOutboundIP() (net.IP, error) {
	for i := 1; i < 100; i++ {
		conn, err := net.Dial("udp", "8.8.8.8:1")
		if err == nil {
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			return localAddr.IP, nil
		}
		defer conn.Close()
	}
	return nil, fmt.Errorf("failed to get current ip address")
}
