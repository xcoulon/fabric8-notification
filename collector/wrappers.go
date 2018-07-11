package collector

import (
	"context"

	"github.com/fabric8-services/fabric8-notification/configuration"
)

func ConfiguredVars(config *configuration.Data, resolver ReceiverResolver) ReceiverResolver {
	return func(ctx context.Context, id string) ([]Receiver, map[string]interface{}, error) {
		r, v, err := resolver(ctx, id)
		if err != nil {
			return r, v, err
		}
		if v == nil {
			v = map[string]interface{}{}
		}
		v["webURL"] = config.GetWebURL()
		return r, v, err
	}
}
