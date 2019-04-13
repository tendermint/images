package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Message struct {
	Channel string `json:"channel"`
	Thread  string `json:"thread_ts,omitempty"`
	Text    string `json:"text"`
}

type SlackResponse struct {
	Ts string `json:"ts"`
}

func main() {
	fail_message := "Simulation with *seed " + os.Args[2] + "* failed! To replicate, run the command:  ```" +
		`go test ./cmd/gaia/app -run TestFullGaiaSimulation \
								-SimulationEnabled=true \
								-SimulationNumBlocks=` + os.Getenv("BLOCKS") + ` \
								-SimulationVerbose=true \
								-SimulationCommit=true \
								-SimulationSeed=` + os.Args[2] + ` \
								-SimulationPeriod=` + os.Getenv("PERIOD") + ` \
								-v -timeout 24h` + "```"

	var message string
	if os.Args[1] == "0" {
		message = "Seed " + os.Args[2] + " *PASS*"
	} else if os.Args[1] == "1" {
		message = fail_message
	} else {
		message = "Host finished running seeds: " + os.Args[2]
	}

	msg_ts := ""
	if len(os.Args) == 4 {
		msg_ts = os.Args[3]
	}
	msg := Message{os.Getenv("SLACK_CHANNEL_ID"), msg_ts, message}

	encoded_payload, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	slack_url := "https://slack.com/api/chat.postMessage"

	u, _ := url.ParseRequestURI(slack_url)
	url_str := u.String()

	req, _ := http.NewRequest("POST", url_str, bytes.NewBuffer(encoded_payload))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SLACK_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data SlackResponse
	json.Unmarshal(body, &data)
	fmt.Printf("%v", data.Ts)
}
