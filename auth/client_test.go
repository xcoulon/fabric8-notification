package auth_test

import (
	"context"
	"testing"

	"github.com/fabric8-services/fabric8-notification/auth"
	"github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	openshiftIOAPI = "http://auth.openshift.io"
)

func createClient(t *testing.T) *api.Client {
	c, err := auth.NewCachedClient(openshiftIOAPI)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSpaceCollaborators(t *testing.T) {

	c := createClient(t)
	id, _ := uuid.FromString("020f756e-b51a-4b43-b113-45cec16b9ce9")

	u, err := auth.GetSpaceCollaborators(context.Background(), c, id)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, len(u.Data) > 10)
}
