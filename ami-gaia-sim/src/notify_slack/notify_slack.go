package main

import (
	"flag"
	"fmt"
	"github.com/nlopes/slack"
)

func main() {
	message_ts := flag.String("ts", "", "Message timestamp")
	channel_id := flag.String("channel-id", "", "Channel ID")
	seed_failed := flag.Bool("failed", false, "Failed seed flag")
	seed_num := flag.String("seed-num", "", "Simulation seed parameter")
	num_blocks := flag.String("num-blocks", "", "Simulation blocks parameter")
	period := flag.String("period", "", "Simulation period parameter")
	slack_token := flag.String("slack-token", "", "Slack token")
	seeds := flag.String("seeds", "", "String representation of the list of seeds")
	log_file_path := flag.String("log-file-path", "", "Path to simulation log file")

	flag.Parse()

	var message string
	if *seed_failed {
		message = "Simulation with *seed " + *seed_num + "* failed! To replicate, run the command:  ```" +
			`go test github.com/cosmos/cosmos-sdk/simapp -run TestFullAppSimulation \
					-SimulationEnabled=true -SimulationNumBlocks=` + *num_blocks + ` \
					-SimulationVerbose=true -SimulationCommit=true \
					-SimulationSeed=` + *seed_num + ` \
					-SimulationPeriod=` + *period + ` \
					-v -timeout 24h` + "```"

	} else {
		message = "Host finished running seeds: " + *seeds
	}

	slack_api := slack.New(*slack_token)

	if len(*message_ts) == 0 {
		_, _, err := slack_api.PostMessage(*channel_id, slack.MsgOptionText(message, false))
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

	} else if *log_file_path != "" {
		log_file_params := slack.FileUploadParameters{
			Title: "Simulation log",
			File:  *log_file_path,
		}
		_, err := slack_api.UploadFile(log_file_params)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	} else {
		_, _, err := slack_api.PostMessage(*channel_id, slack.MsgOptionText(message, false), slack.MsgOptionTS(*message_ts))
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}
}
