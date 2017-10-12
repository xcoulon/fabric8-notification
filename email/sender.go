package email

import (
	"context"

	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-wit/log"
	"github.com/mattbaird/gochimp"
)

type Sender interface {
	Send(ctx context.Context, subject string, body string, headers map[string]string, receivers []collector.Receiver)
}

func NewMandrillSender(apiKey string) (Sender, error) {
	api, err := gochimp.NewMandrill(apiKey)
	if err != nil {
		return nil, err
	}
	return &MandrillSender{mandrillAPI: api}, nil
}

type MandrillSender struct {
	mandrillAPI *gochimp.MandrillAPI
}

func (m *MandrillSender) Send(ctx context.Context, subject string, body string, headers map[string]string, receivers []collector.Receiver) {
	recipients := toRecipients(receivers)

	message := gochimp.Message{
		Html:      body,
		Subject:   subject,
		FromEmail: "noreply@notify.openshift.io",
		FromName:  "openshift.io",
		To:        recipients,
		Headers:   headers,
	}

	resps, err := m.mandrillAPI.MessageSend(message, false)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "error sending messages")
	}
	for _, resp := range resps {
		if resp.Status != "sent" && resp.Status != "queued" {
			log.Error(ctx, map[string]interface{}{
				"recipient":    resp.Email,
				"recipient_id": resp.Id,
				"status":       resp.Status,
				"rejected":     resp.RejectedReason,
			}, "sent message failed")

		} else {
			log.Info(ctx, map[string]interface{}{
				"recipient":    resp.Email,
				"recipient_id": resp.Id,
				"status":       resp.Status,
			}, "sent message")
		}
	}
}

func toRecipients(receivers []collector.Receiver) []gochimp.Recipient {
	var recipients []gochimp.Recipient
	for _, r := range receivers {
		recipients = append(recipients, gochimp.Recipient{Name: r.FullName, Email: r.EMail})
	}
	return recipients
}
