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
	numSeeds                  = 900
	numJobs                   = numSeeds
	instanceShutdownBehaviour = "terminate"
)

var err error

func makeSeedLists() map[int]string {
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
	var (
		amiID         string
		channelID     string
		slackToken    string
		numBlocks     string
		simPeriod     string
		gitRevision   string
		messageTS     string
		logObjPrefix  string
		instanceIndex []int
		notifyOnly    bool
		genesis       bool
		sessionEC2    = ec2.New(session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})))
		ec2Instances *ec2.Reservation
	)

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

		// DO NOT REMOVE. Using this output to set an environment variable
		// Env variable will be used by subsequent runs of startsim
		fmt.Println(msgTS)
		os.Exit(0)
	}

	messageTS = os.Getenv("MSGTS")
	if _, err = slackMessage(slackToken, channelID, &messageTS, "Spinning up simulation environments!"); err != nil {
		log.Fatal("Could not report back to slack: " + err.Error())
	}

	if amiID, err = getAmiId(gitRevision, sessionEC2); err != nil {
		log.Fatal(err.Error())
	}

	logObjPrefix = fmt.Sprintf("simID-%s", os.Getenv("CIRCLE_BUILD_NUM"))
	seedLists := makeSeedLists()
	for index := range seedLists {
		var userData strings.Builder
		userData.WriteString("#!/bin/bash \n")
		userData.WriteString("cd /home/ec2-user/go/src/github.com/cosmos/cosmos-sdk \n")
		userData.WriteString("source /etc/profile.d/set_env.sh \n")
		userData.WriteString(buildCommand(numJobs, logObjPrefix, seedLists[index], slackToken, channelID, messageTS, numBlocks, simPeriod, genesis))
		userData.WriteString("shutdown -h now")

		config := &ec2.RunInstancesInput{
			InstanceInitiatedShutdownBehavior: aws.String(instanceShutdownBehaviour),
			InstanceType:                      aws.String("c4.8xlarge"),
			ImageId:                           aws.String(amiID),
			KeyName:                           aws.String("wallet-nodes"),
			MaxCount:                          aws.Int64(1),
			MinCount:                          aws.Int64(1),
			UserData:                          aws.String(base64.StdEncoding.EncodeToString([]byte(userData.String()))),
			IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
				Name: aws.String("gaia-simulation"),
			}}

		if ec2Instances, err = sessionEC2.RunInstances(config); err != nil {
			// Checking aws error code to see if we have reached the EC2 instance limit for this instance type
			// Crashing out of the program is not desirable. We can run simulation with a lower number of seeds
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == "InstanceLimitExceeded" {
					log.Println(awsErr.Error())
					break
				} else {
					log.Fatal(awsErr.Error())
				}
			} else {
				log.Fatal(err.Error())
			}
		}
		for i := range ec2Instances.Instances {
			log.Println(*ec2Instances.Instances[i].InstanceId)
		}
		instanceIndex = append(instanceIndex, index)
	}
	if len(instanceIndex) > 1 {
		sendSqsMsg(instanceIndex)
	}
}
