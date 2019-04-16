package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type SimConfig struct {
	SlackChannelId string
	SlackToken     string
	Blocks         string
	Period         string
	MessageTs      string
	GitRevision    string
}

type SlackPayload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
	Thread  string `json:"thread_ts,omitempty"`
}

type SlackResponse struct {
	Ts string `json:"ts"`
}

const num_seeds = 36

func make_ranges() map[int]string {
	machines := make(map[int]string)
	var str strings.Builder
	index := 0
	for i := 0; i < num_seeds; i++ {
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

func push_to_slack(slack_token string, slack_channel_id string, message string, message_ts string) (string, error) {
	slack_url := "https://slack.com/api/chat.postMessage"

	slack_payload, json_err := json.Marshal(SlackPayload{
		Channel: slack_channel_id,
		Text:    message,
	})
	if json_err != nil {
		return "", json_err
	}

	request, _ := http.NewRequest("POST", slack_url, bytes.NewBuffer(slack_payload))
	request.Header.Set("Authorization", "Bearer "+slack_token)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	http_client := &http.Client{Timeout: 10 * time.Second}
	resp, err := http_client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var data SlackResponse
	json_err = json.NewDecoder(resp.Body).Decode(&data)
	if json_err != nil {
		return "", json_err
	}
	return data.Ts, nil
}

func get_ami_id(git_revision string, svc *ec2.EC2) (string, error) {
	input := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("name"),
				Values: []*string{
					aws.String("gaia-sim-" + git_revision),
				},
			},
		},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		return "", err
	}
	return *result.Images[0].ImageId, nil
}

func main() {

	// Yes I know the Args usage is fugly. I will make it nice and pretty
	// with something like cobra, but this needs to get out the door and be
	// useful now.
	// It's safe from shennanigans as the program will only be called from Circle
	// fixed configuration.
	msg_ts, slack_err := push_to_slack(os.Args[5], os.Args[4], "Spinning up simulation environments!", "")
	if slack_err != nil {
		fmt.Println("Could not report back to slack: " + slack_err.Error())
		os.Exit(1)
	}

	conf := new(SimConfig)
	conf.Blocks = os.Args[1]
	conf.Period = os.Args[2]
	conf.GitRevision = os.Args[3]
	conf.SlackChannelId = os.Args[4]
	conf.SlackToken = os.Args[5]
	conf.MessageTs = msg_ts

	svc := ec2.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})))
	
	ami_id, err := get_ami_id(conf.GitRevision, svc)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	seeds := make_ranges()
	for rng := range seeds {
		var usr_data strings.Builder
		usr_data.WriteString("#!/bin/bash\n")
		usr_data.WriteString("export BLOCKS=" + conf.Blocks + "\n")
		usr_data.WriteString("export PERIOD=" + conf.Period + "\n")
		usr_data.WriteString("export SLACK_TOKEN=" + conf.SlackToken + "\n")
		usr_data.WriteString("export SLACK_CHANNEL_ID=" + conf.SlackChannelId + "\n")
		usr_data.WriteString("export SLACK_MSG_TS=" + conf.MessageTs + "\n")
		usr_data.WriteString("cd /home/ec2-user/go/src/github.com/cosmos/cosmos-sdk/\n")
		usr_data.WriteString("./multisim.sh " + seeds[rng] + "> /home/ec2-user/sim_out 2>&1")

		fmt.Println(usr_data.String())
		config := &ec2.RunInstancesInput{
			InstanceInitiatedShutdownBehavior: aws.String("stop"),
			InstanceType:                      aws.String("c4.8xlarge"),
			ImageId:                           aws.String(ami_id),
			KeyName:                           aws.String("wallet-nodes"),
			MaxCount:                          aws.Int64(1),
			MinCount:                          aws.Int64(1),
			UserData:                          aws.String(base64.StdEncoding.EncodeToString([]byte(usr_data.String()))),
		}
		result, err := svc.RunInstances(config)
		if err != nil {
			fmt.Println(err.Error())
		}

		for i := range result.Instances {
			fmt.Println(*result.Instances[i].InstanceId)
		}
	}
}
