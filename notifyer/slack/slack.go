// Package slack implements slack
package slack

import (
	"atlassian_backup/lib/utils"
	"errors"
	"os"

	"github.com/slack-go/slack"
)

const slackSender = "Atlassian backup service"

type Notifyer struct {
	webhookUrl string
}

func New() (*Notifyer, error) {
	webhook, ok := os.LookupEnv("SLACK_WEBHOOK_URL")
	if !ok {
		return nil, errors.New("Slack webhook URL is not specified")
	}

	return &Notifyer{
		webhookUrl: webhook,
	}, nil
}

func (n *Notifyer) Send(text string) error {
	message := &slack.WebhookMessage{
		Username: slackSender,
		Text:     text,
	}

	err := slack.PostWebhook(n.webhookUrl, message)
	if err != nil {
		return utils.Wrap("can't send notification to Slack", err)
	}
	return nil
}
