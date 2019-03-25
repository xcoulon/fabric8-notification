package validator

import (
	"context"
	"github.com/fabric8-services/fabric8-common/errors"
)

func ValidateUser(context context.Context, custom map[string]interface{}) error {
	_, found := custom["verifyURL"]
	if !found {
		return errors.NewBadParameterError("data.attributes.custom.verifyURL", "nil")
	}
	return nil
}
