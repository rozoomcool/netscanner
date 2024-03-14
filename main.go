package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {

	fmt.Println("Start scan network")
	fmt.Println("Let's scan your local network!!!")

	scanIP("localhost")

}

func scanIP(host string) {
	ports := make(chan int)
	func() {
		defer close(ports)
		checkAddressPorts(host, ports)
	}()

	for p := range ports {
		if p != 0 {
			fmt.Printf("PORT %v is open\n", p)
		}
	}
}

func scanLocalNetwork() {

	hosts := getEnabledHosts()

	for _, host := range hosts {
		if host == "" {
			fmt.Println(host)
			continue
		}

		ports := make(chan int)

		func() {
			defer close(ports)
			checkAddressPorts(host, ports)
		}()

		fmt.Printf("HOST:: %v\n", host)
		for i := range ports {
			if i != 0 {
				fmt.Printf("\tport %d is open\n", i)
			}
		}
	}
}

func getEnabledHosts() []string {
	baseIP := getPrefIP()

	var hosts []string

	for i := 1; i <= 255; i++ {
		cmd := exec.Command("ping", "-c", "1", "-W", "1", fmt.Sprintf("%v:%d", baseIP, i))
		if err := cmd.Run(); err != nil {
			continue
		}

		if err := cmd.Wait(); err != nil {
			continue
		}

		hosts = append(hosts, fmt.Sprintf("%v:%d", baseIP, i))
	}

	return hosts
}

func getPrefIP() string {
	addr, err := GetOutboundIP()
	if err != nil {
		fmt.Println(err)
	}
	ip := addr.String()

	ipArr := strings.Split(ip, ".")
	baseIP := strings.Join(ipArr[0:len(ipArr)-1], ".") + "."

	return baseIP
}

func checkAddressPorts(host string, result chan int) {
	for i := 1; i <= 65536; i++ {
		go func(i int) {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%d", host, i), time.Second*3)
			if err != nil {
			} else {
				conn.LocalAddr()
				result <- i
			}
		}(i)
	}
}

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP, nil
	}
	defer conn.Close()
	return nil, fmt.Errorf("failed to get current ip address")
}
