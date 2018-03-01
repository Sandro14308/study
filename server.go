//nothing
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type PageVariables struct {
	PageTitle string
	Answer    string
}

func main() {
	fmt.Println("Server start 8080")
	http.HandleFunc("/", HomePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	Title := "Lol kek"
	MyPageVariables := PageVariables{
		PageTitle: Title,
	}
	t, err := template.ParseFiles("index.html") //parse the html file homepage.html
	if err != nil {                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}
