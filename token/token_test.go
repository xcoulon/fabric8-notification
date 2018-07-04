package token_test

import (
	"os"
	"testing"

	"github.com/fabric8-services/fabric8-notification/configuration"
	"github.com/fabric8-services/fabric8-notification/token"
	"github.com/stretchr/testify/require"
)

const (
	OpenShiftIOAuthAPI = "https://auth.openshift.io"
)

func TestManager(t *testing.T) {

	old := os.Getenv("F8_AUTH_URL")
	defer os.Setenv("F8_AUTH_URL", old)

	os.Setenv("F8_AUTH_URL", OpenShiftIOAuthAPI)
	config, err := configuration.GetData()
	require.NoError(t, err)
	require.NotNil(t, config)

	manager, err := token.NewManager(config)
	require.NoError(t, err)

	keys := manager.PublicKeys()
	require.NotEmpty(t, keys)
	for _, k := range keys {
		require.NotNil(t, k.N)
	}
}
