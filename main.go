package main

// DON'T FORGET TO CHECK https://www.google.com/settings/security/lesssecureapps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

type site struct {
	Mail_reception string   `json:"mail_reception"`
	Mail_envoie    string   `json:"mail_envoie"`
	Mdp_envoie     string   `json:"mdp_envoie"`
	Site           []string `json:"site"`
}

func main() {
	doEvery(10 * time.Second)

}

func doEvery(d time.Duration) {
	for x := range time.Tick(d) {
		checkSite(x)
	}
}

func checkSite(t time.Time) {
	sites := getJson("./config.json")

	for i := 0; i < len(sites.Site); i++ {
		status := pingSite(sites.Site[i])

		if status {
			send(sites.Site[i]+" est Joignable", sites)
		} else {
			send(sites.Site[i]+"n'est plus joignable", sites)
		}
	}
}

func pingSite(url string) bool {
	var bodyString bool
	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode == 200 { // OK
		bodyString = true

	} else {
		bodyString = false
	}
	return bodyString

}

func getJson(url string) site {
	file, err := ioutil.ReadFile(url)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var jsontype site
	json.Unmarshal(file, &jsontype)
	return jsontype
}

func send(body string, config site) {
	from := config.Mail_envoie
	pass := config.Mdp_envoie
	to := config.Mail_reception

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit http://foobarbazz.mailinator.com")
}
