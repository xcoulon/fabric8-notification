package collector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fabric8-services/fabric8-notification/configuration"
)

func TestConfigureVarsSetWebURL(t *testing.T) {
	config, err := configuration.NewData()
	if err != nil {
		assert.NoError(t, err)
	}
	u, v, err := ConfiguredVars(config, EmptyResolver)(context.Background(), "id")
	assert.NoError(t, err)
	assert.Len(t, u, 0)
	assert.Len(t, v, 1)
	assert.NotEmpty(t, v["webURL"])
}

func TestConfigureVarsNilResolverVars(t *testing.T) {
	config, err := configuration.NewData()
	if err != nil {
		assert.NoError(t, err)
	}
	u, v, err := ConfiguredVars(config, NilVarResolver)(context.Background(), "id")
	assert.NoError(t, err)
	assert.Len(t, u, 0)
	assert.Len(t, v, 1)
	assert.NotEmpty(t, v["webURL"])
}

func EmptyResolver(context.Context, string) (users []Receiver, templateValues map[string]interface{}, err error) {
	return []Receiver{}, map[string]interface{}{}, nil
}

func NilVarResolver(context.Context, string) (users []Receiver, templateValues map[string]interface{}, err error) {
	return []Receiver{}, nil, nil
}
