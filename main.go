package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"runtime"
	"strings"
	"text/template"
)

// Resp xxx
type Resp struct {
	Answer AnswerResp `json:"answer"`
}

// AnswerResp xxx
type AnswerResp struct {
	Class        string      `json:"__class"`
	Availability []AvailResp `json:"availability"`
}

// AvailResp xxx
type AvailResp struct {
	Reference string     `json:"reference"`
	MetaZones []ZoneResp `json:"metaZones"`
	Zones     []ZoneResp `json:"zones"`
}

// ZoneResp xxx
type ZoneResp struct {
	Availability string `json:"availability"`
	Zone         string `json:"zone"`
}

func main() {
	URL := "https://ws.ovh.com/dedicated/r2/ws.dispatcher/getAvailability2"
	SEND := "https://www.kimsufi.com/fr/commande/kimsufi.xml?reference=150sk30&quantity=1"
	fmt.Println(URL)

	// Create HTTP request
	client := http.Client{}
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Execute request
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var m Resp
	err = json.Unmarshal(body, &m)

	for _, server := range m.Answer.Availability {
		if server.Reference == "150sk30" {
			fmt.Println(server.MetaZones)
			for _, zone := range server.Zones {
				if (zone.Zone == "fr" || zone.Zone == "westernEurope") && zone.Availability != "unknown" {
					// send dispo
					SendEmail("localhost", 25, "plex", "plex190058", []string{"bpetetot@gmail.com"}, "Kimsufi Dispo", SEND)
				}
			}
			for _, metazone := range server.MetaZones {
				if (metazone.Zone == "fr" || metazone.Zone == "westernEurope") && metazone.Availability != "unknown" {
					// send dispo
					SendEmail("localhost", 25, "plex", "plex190058", []string{"bpetetot@gmail.com"}, "Kimsufi Dispo", SEND)
				}
			}
		}
	}

	SendEmail("localhost", 25, "plex", "plex190058", []string{"bpetetot@gmail.com"}, "Kimsufi Dispo", SEND)
}

func catchPanic(err *error, functionName string) {
	if r := recover(); r != nil {
		fmt.Printf("%s : PANIC Defered : %v\n", functionName, r)

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		fmt.Printf("%s : ERROR : %v\n", functionName, *err)

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))
	}
}

// SendEmail xxx
func SendEmail(host string, port int, userName string, password string, to []string, subject string, message string) (err error) {
	defer catchPanic(&err, "SendEmail")

	parameters := struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		userName,
		strings.Join([]string(to), ","),
		subject,
		message,
	}

	buffer := new(bytes.Buffer)

	template := template.Must(template.New("emailTemplate").Parse(emailScript()))
	template.Execute(buffer, &parameters)

	auth := smtp.PlainAuth("", userName, password, host)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		userName,
		to,
		buffer.Bytes())

	return err
}

func emailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}`
}
