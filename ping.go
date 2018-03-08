package main

import (
	"fmt"
	"os/exec"
	"strings"
	"strconv"
	"time"

)

func ping(i int, msg chan<- string ) {
	tmp := strconv.Itoa(i)
	
	host :="192.168.0." + tmp
	
	out, _ := exec.Command("ping", host, "-c 1", "-i 1", "-w 10").Output()
	
	if strings.Contains(string(out), "Destination Host Unreachable") {
	 	  msg <- "HOST: " + host + " NONE"
	}else {
		
	    msg <- "HOST: " + host + " available"

	}
	return
}

func sleep(seconds int, endSignal chan<- bool) {
    time.Sleep(time.Duration(seconds) * time.Second)
    endSignal <- true
}

func main() {
	fmt.Println("Starting program")

	


/*      first part sleep



	endSignal := make(chan bool, 1)
    go sleep(3, endSignal)
    var end bool
*/	


    msg := make(chan string)

	for i := 100; i < 151; i++ {
		go ping(i, msg)
	}	



	for i := 0; i < 150; i++ {
		fmt.Println(<-msg)
	}


/*      second part sleep




    for !end {
        select {
        case end = <-endSignal:
            fmt.Println("The end!")
        case <-time.After(5 * time.Second):
            fmt.Println("There's no more time to this. Exiting!")
            end = true
        }
    }*/



	//select{}
}
	
