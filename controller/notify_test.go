package controller

import (
	"context"
	"fmt"
	"testing"

	"github.com/fabric8-services/fabric8-notification/app"
	"github.com/fabric8-services/fabric8-notification/app/test"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/email"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/goadesign/goa"
)

func TestNotifySendUnknownType(t *testing.T) {

	resolvers := &collector.LocalRegistry{}
	typeRegistry := &template.AssetRegistry{}
	notifier := &email.CallbackNotifier{Callback: func(ctx context.Context, notification email.Notification) { fmt.Println(notification) }}
	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13131",
				Type: "unknown.create",
			},
			Type: "notifications",
		},
	}

	test.SendNotifyBadRequest(t, nil, nil, ctrl, payload)
}
