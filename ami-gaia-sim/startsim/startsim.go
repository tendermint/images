package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"


	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// If the number of jobs is < the number of seeds, simulation will crash
	numSeeds                  = 384
	numJobs                   = numSeeds
	instanceShutdownBehaviour = "terminate"
)

var (
	channelID    string
	slackToken   string
	numBlocks    string
	simPeriod    string
	gitRevision  string
	messageTS    string
	logObjPrefix string
	err          error
	notifyOnly   bool
	genesis      bool
)

func makeRanges() map[int]string {
	machines := make(map[int]string)
	var str strings.Builder
	index := 0
	for i := 0; i <= numSeeds; i++ {
		if i != 0 && math.Mod(float64(i), 35) == 0 {
			machines[index] = strings.TrimRight(str.String(), ",")
			str.Reset()
			index++
		}
		str.WriteString(strconv.Itoa(i) + ",")
	}

	if str.String() != "" {
		machines[index] = strings.TrimRight(str.String(), ",")
	}
	return machines
}

func getAmiId(gitRevision string, svc *ec2.EC2) (string, error) {
	var imageID *ec2.DescribeImagesOutput

	input := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("name"),
				Values: []*string{
					aws.String("gaia-sim-" + gitRevision),
				},
			},
		},
	}
	if imageID, err = svc.DescribeImages(input); err != nil {
		return "", err
	}
	return *imageID.Images[0].ImageId, nil
}

func buildCommand(jobs int, logObjKey, seeds, token, channel, timeStamp, blocks, period string, genesis bool) string {
	if genesis {
		return fmt.Sprintf("runsim -log \"%s\" -j %d -seeds \"%s\" -g /home/ec2-user/genesis.json "+
			"-slack \"%s,%s,%s\" github.com/cosmos/cosmos-sdk/simapp %s %s TestFullAppSimulation;",
			logObjKey, jobs, seeds, token, channel, timeStamp, blocks, period)
	}
	return fmt.Sprintf("runsim -log \"%s\" -j %d -seeds \"%s\" -slack \"%s,%s,%s\" github.com/cosmos/cosmos-sdk/simapp %s %s TestFullAppSimulation;",
		logObjKey, jobs, seeds, token, channel, timeStamp, blocks, period)
}

func main() {
	flag.StringVar(&slackToken, "s", "", "Slack token")
	flag.StringVar(&channelID, "c", "", "Slack channel ID")
	flag.StringVar(&numBlocks, "b", "", "Number of blocks to simulate")
	flag.StringVar(&simPeriod, "p", "", "Simulation invariant check period")
	flag.StringVar(&gitRevision, "g", "", "The git revision on which the simulation is run")
	flag.BoolVar(&notifyOnly, "notify", false, "Send notification and exit")
	flag.BoolVar(&genesis, "gen", false, "Use genesis file in simulation")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(),
			`Usage: %s [-notify] [-gen] [-s slacktoken] [-c channelID] [-b numblocks] [-p simperiod] [-g gitrevision]`, filepath.Base(os.Args[0]))
	}
	flag.Parse()

	if notifyOnly {
		msgTS, err := slackMessage(slackToken, channelID, nil,
			fmt.Sprintf("*Starting simulation #%s.* SDK hash/tag/branch: `%s`. <%s|Circle build url>\nblocks:\t`%s`\nperiod:\t`%s`\nseeds:\t`%d`",
				os.Getenv("CIRCLE_BUILD_NUM"), gitRevision, os.Getenv("CIRCLE_BUILD_URL"), numBlocks, simPeriod, numSeeds))

		if err != nil {
			log.Printf("ERROR: sending slack message: %v", err)
		}

		// DO NOT REMOVE. This output is used by other tools that run after this one.
		fmt.Println(msgTS)
		os.Exit(0)
	}

	messageTS = os.Getenv("MSGTS")
	_, err = slackMessage(slackToken, channelID, &messageTS, "Spinning up simulation environments!")
	if err != nil {
		log.Fatal("Could not report back to slack: " + err.Error())
	}

	svc := ec2.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})))

	amiId, err := getAmiId(gitRevision, svc)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Fatal(awsErr.Error())
		}
		log.Fatal(err.Error())
	}

	logObjPrefix = fmt.Sprintf("simID-%s", os.Getenv("CIRCLE_BUILD_NUM"))
	seeds := makeRanges()
	for rng := range seeds {
		var userData strings.Builder
		userData.WriteString("#!/bin/bash \n")
		userData.WriteString("cd /home/ec2-user/go/src/github.com/cosmos/cosmos-sdk \n")
		userData.WriteString("source /etc/profile.d/set_env.sh \n")
		userData.WriteString(buildCommand(numJobs, logObjPrefix, seeds[rng], slackToken, channelID, messageTS, numBlocks, simPeriod, genesis))
		userData.WriteString("shutdown -h now")

		config := &ec2.RunInstancesInput{
			InstanceInitiatedShutdownBehavior: aws.String(instanceShutdownBehaviour),
			InstanceType:                      aws.String("c4.8xlarge"),
			ImageId:                           aws.String(amiId),
			KeyName:                           aws.String("wallet-nodes"),
			MaxCount:                          aws.Int64(1),
			MinCount:                          aws.Int64(1),
			UserData:                          aws.String(base64.StdEncoding.EncodeToString([]byte(userData.String()))),

			IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
				Name: aws.String("gaia-simulation"),
			}}
		result, err := svc.RunInstances(config)
		if err != nil {
			log.Fatal(err.Error())
		}

		for i := range result.Instances {
			log.Println(*result.Instances[i].InstanceId)
		}
	}
	if len(seeds) > 1 {
		sendSqsMsg(seeds)
	}
}
