package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
					SendEmail("bpetetot@gmail.com", SEND)
				}
			}
			for _, metazone := range server.MetaZones {
				if (metazone.Zone == "fr" || metazone.Zone == "westernEurope") && metazone.Availability != "unknown" {
					// send dispo
					SendEmail("bpetetot@gmail.com", SEND)
				}
			}
		}
	}

	SendEmail("bpetetot@gmail.com", SEND)
}

// SendEmail xxx
func SendEmail(to string, message string) {
	if err := exec.Command("echo \"" + message + "\" | mail -s 'Dispo' -aFrom:Kimsufi Available " + to).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Successfully halved image in size")
}
