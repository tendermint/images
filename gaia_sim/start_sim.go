package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type SlackPayload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func make_ranges() map[int]string {
	machines := make(map[int]string)
	var str strings.Builder
	index := 0
	for i := 0; i < 36; i++ {
		if math.Mod(float64(i), 35) == 0 {
			if index != 0 {
				machines[index] = str.String()
			}
			str.Reset()
			index++
		}
		str.WriteString(strconv.Itoa(i) + " ")
	}

	return machines
}

func push_to_slack(slack_token string, slack_channel_id string) {
	slack_url := "https://slack.com/api/chat.postMessage"

	json_payload, _ := json.Marshal(
		SlackPayload{
			Channel: slack_channel_id,
			Text:    "Simulation started!",
		})

	slack_payload := bytes.NewBuffer(json_payload)

	request, _ := http.NewRequest("POST", slack_url, slack_payload)
	request.Header.Set("Authorization", "Bearer "+slack_token)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	http_client := &http.Client{Timeout: 10 * time.Second}
	response, err := http_client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()
}

func start_sim(slack_token string, slack_channel_id string) {
	ami_id, has_ami := os.LookupEnv("AMI")

	// TODO: decide how logging will work
	if !has_ami {
		fmt.Println("FMT: No AMI!")
		log.Println("LOG: No AMI!")
		os.Exit(1)
	}

	blocks, has_blocks := os.LookupEnv("BLOCKS")
	period, has_period := os.LookupEnv("PERIOD")

	if !has_blocks || !has_period {
		blocks = "100"
		period = "10"
	}

	seeds := make_ranges()
	for rng := range seeds {

		user_data := `#!/bin/bash
                  export BLOCKS=` + blocks + `
                  export PERIOD=` + period + `
                  export SEEDS=(` + seeds[rng] + `)
                  export SLACK_TOKEN=` + slack_token + `
                  export SLACK_CHANNEL_ID=` + slack_channel_id + `
                  /home/ec2-user/go/src/github.com/cosmos/cosmos-sdk/multisim.sh > /home/ec2-user/sim_out 2>&1`

		svc := ec2.New(session.Must(session.NewSession()))
		config := &ec2.RunInstancesInput{
			InstanceInitiatedShutdownBehavior: aws.String("terminate"),
			InstanceType:                      aws.String("c4.8xlarge"),
			ImageId:                           aws.String(ami_id),
			KeyName:                           aws.String("wallet-nodes"),
			MaxCount:                          aws.Int64(1),
			MinCount:                          aws.Int64(1),
			UserData:                          aws.String(base64.StdEncoding.EncodeToString([]byte(user_data))),
		}

		result, _ := svc.RunInstances(config)

		for i := range result.Instances {
			fmt.Println(result.Instances[i].InstanceId)
		}
	}
}

func main() {
	slack_token, _ := os.LookupEnv("SLACK_TOKEN")
	slack_channel_id, _ := os.LookupEnv("SLACK_CHANNEL_ID")

	push_to_slack(slack_token, slack_channel_id)
	start_sim(slack_token, slack_channel_id)
}
