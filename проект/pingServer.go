package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type PageVariables struct {
	PageTitle string
	IpAddress string
}

func startPage(w http.ResponseWriter, r *http.Request) {
	Title := "Find hosts"
	ip := "192.168.0.1"
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

func main() {
	fmt.Println("Server start. Port *:8080")
	http.HandleFunc("/", startPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
