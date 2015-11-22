package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	//SEND := "https://www.kimsufi.com/fr/commande/kimsufi.xml?reference=150sk30&quantity=1"
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
					SendEmail("bpetetot@gmail.com")
				}
			}
			for _, metazone := range server.MetaZones {
				if (metazone.Zone == "fr" || metazone.Zone == "westernEurope") && metazone.Availability != "unknown" {
					// send dispo
					SendEmail("bpetetot@gmail.com")
				}
			}
		}
	}

	SendEmail("bpetetot@gmail.com")
}

// SendEmail xxx
func SendEmail(to string) {

	echo := exec.Command("echo", "Test")
	mail := exec.Command("mail", "-s", "Kimsufi", to)
	output, err := pipeCommands(echo, mail)

	if err != nil {
		println("error : ")
		println(err)
	} else {
		println("output : ")
		print(string(output))
	}

}

func pipeCommands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out

		defer command.Wait() // Doesn't block
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil, err
	}
	return final, nil
}
