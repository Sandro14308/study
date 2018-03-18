package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	fastping "github.com/tatsushid/go-fastping"
)

type PageVariables struct {
	PageTitle string
	IpAddress string
}

func logger() {

}

func startPage(w http.ResponseWriter, r *http.Request) {
	Title := "Find hosts"
	ip := "192.168.1.0"
	StartPageVariable := PageVariables{
		PageTitle: Title,
		IpAddress: ip,
	}

	t, err := template.ParseFiles("page.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, StartPageVariable)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

/*
func ProcessingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("ip:", r.Form["ipAddress"])
	}
	ipAddress := strings.Join(r.Form["ipAddress"], "")

	var wg sync.WaitGroup
	waitGroupLength := 255
	wg.Add(waitGroupLength)

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

}*/

func Ping(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	result := make(chan string, 2)
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", string(body[:]))
	if err != nil {
		result <- "ERR"
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		result <- "OK"
	}
	p.OnIdle = func() {
		result <- "ERR"
	}
	err = p.Run()
	if err != nil {
		result <- "ERR"
	}
	res := <-result
	fmt.Println(string(body[:]), res)
	io.WriteString(w, res)
	log.Println(string(body[:]), " ", res)
}

func main() {
	f, err := os.OpenFile("log2.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)

	fmt.Println("Server start. Port *:8080")
	http.HandleFunc("/", startPage)
	//http.HandleFunc("/processing", ProcessingPage)
	http.HandleFunc("/ping", Ping)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
