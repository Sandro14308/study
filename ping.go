package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	fastping "github.com/tatsushid/go-fastping"
)

func main() {

	var wg sync.WaitGroup
	waitGroupLength := 8
	wg.Add(waitGroupLength)

	ipAddress := os.Args[1]
	if strings.HasSuffix(ipAddress, ".0") {

		ipAddress = strings.TrimSuffix(ipAddress, "0")

		result := make(chan string, 1)
		errChannel := make(chan error, 1)

		for i := 1; i < 255; i++ {
			tmp := ipAddress + strconv.Itoa(i)
			go func(tmp string, i int) {
				time.Sleep(time.Duration(waitGroupLength - i - 99))
				time.Sleep(0)
				p := fastping.NewPinger()
				ra, err := net.ResolveIPAddr("ip4:icmp", tmp)
				if err != nil {
					errChannel <- err
				}
				p.AddIPAddr(ra)
				p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
					result <- "IP Addr: " + tmp
				}
				p.OnIdle = func() {
					result <- "Ip Addr: " + tmp + " NONE"
				}
				err = p.Run()
				if err != nil {
					errChannel <- err
				}
				wg.Done()
			}(tmp, i)
		}

		go func() {
			wg.Wait()
			close(result)
		}()

		for i := 0; i < 255; i++ {
			select {
			case res := <-result:

				fmt.Println(res)

			case err := <-errChannel:
				if err != nil {
					fmt.Println("error ", err)
					return
				}
			}
		}

	} else {
		fmt.Println("укажите сеть, а не ip адресс")
	}

}
