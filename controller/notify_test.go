package controller

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/fabric8-services/fabric8-notification/app"
	"github.com/fabric8-services/fabric8-notification/app/test"
	"github.com/fabric8-services/fabric8-notification/auth"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/configuration"
	"github.com/fabric8-services/fabric8-notification/email"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/fabric8-services/fabric8-notification/validator"
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

func TestNotifySendWithCustomParam(t *testing.T) {

	config, _ := configuration.GetData()

	resolvers := &collector.LocalRegistry{}
	typeRegistry := &template.AssetRegistry{}
	authClient, _ := auth.NewCachedClient(config.GetWITURL())

	resolvers.Register("user.email.update", collector.ConfiguredVars(config, collector.NewUserResolver(authClient)), validator.ValidateUser)
	notifier := &email.CallbackNotifier{
		Callback: func(ctx context.Context, notification email.Notification) {
			require.NotNil(t, notification.CustomAttributes)
			assert.Equal(t, notification.CustomAttributes["verifyURL"], "https://someurl.openshift.io")
		}}

	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13132",
				Type: "user.email.update",
				Custom: map[string]interface{}{
					"verifyURL": "https://someurl.openshift.io",
				},
			},
			Type: "notifications",
		},
	}

	test.SendNotifyAccepted(t, nil, nil, ctrl, payload)
}

func TestNotifySendWithoutustomParamBadRequest(t *testing.T) {

	config, _ := configuration.GetData()

	resolvers := &collector.LocalRegistry{}
	typeRegistry := &template.AssetRegistry{}
	authClient, _ := auth.NewCachedClient(config.GetWITURL())

	resolvers.Register("user.email.update", collector.ConfiguredVars(config, collector.NewUserResolver(authClient)), validator.ValidateUser)
	notifier := &email.CallbackNotifier{
		Callback: func(ctx context.Context, notification email.Notification) {
			require.NotNil(t, notification.CustomAttributes)
			assert.Equal(t, notification.CustomAttributes["verifyURL"], "https://someurl.openshift.io")
		}}

	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13132",
				Type: "user.email.update",
			},
			Type: "notifications",
		},
	}

	test.SendNotifyBadRequest(t, nil, nil, ctrl, payload)
}

func TestNotifySendWithCustomParamBadRequest(t *testing.T) {

	config, _ := configuration.GetData()

	resolvers := &collector.LocalRegistry{}
	typeRegistry := &template.AssetRegistry{}
	authClient, _ := auth.NewCachedClient(config.GetWITURL())

	resolvers.Register("user.email.update", collector.ConfiguredVars(config, collector.NewUserResolver(authClient)), validator.ValidateUser)
	notifier := &email.CallbackNotifier{
		Callback: func(ctx context.Context, notification email.Notification) {
			require.NotNil(t, notification.CustomAttributes)
			assert.Equal(t, notification.CustomAttributes["verifyURL"], "https://someurl.openshift.io")
		}}

	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13132",
				Type: "user.email.update",
				Custom: map[string]interface{}{
					"somthing_else": "https://someurl.openshift.io",
				},
			},
			Type: "notifications",
		},
	}

	test.SendNotifyBadRequest(t, nil, nil, ctrl, payload)
}

func TestNotifySendWithoutCustomParamSuccess(t *testing.T) {

	config, _ := configuration.GetData()

	resolvers := &collector.LocalRegistry{}
	typeRegistry := &template.AssetRegistry{}
	authClient, _ := auth.NewCachedClient(config.GetWITURL())

	resolvers.Register("workitem.update", collector.ConfiguredVars(config, collector.NewUserResolver(authClient)), nil)
	notifier := &email.CallbackNotifier{
		Callback: func(ctx context.Context, notification email.Notification) {
		}}

	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13132",
				Type: "workitem.update",
			},
			Type: "notifications",
		},
	}

	test.SendNotifyAccepted(t, nil, nil, ctrl, payload)
}
