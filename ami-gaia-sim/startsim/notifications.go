package main

import "github.com/nlopes/slack"

func slackMessage(token string, channel string, threadTS *string, message string) (string, error) {
	client := slack.New(token)
	if threadTS != nil {
		_, respTS, err := client.PostMessage(channel, slack.MsgOptionText(message, false), slack.MsgOptionTS(*threadTS))
		if err != nil {
			return "", err
		}
		return respTS, nil
	} else {
		_, respTS, err := client.PostMessage(channel, slack.MsgOptionText(message, false))
		if err != nil {
			return "", err
		}
		return respTS, nil
	}
}
