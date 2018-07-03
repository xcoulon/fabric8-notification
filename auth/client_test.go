package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/fabric8-services/fabric8-notification/testsupport"
	"github.com/stretchr/testify/require"

	"net/http"
	"testing"

	"github.com/fabric8-services/fabric8-notification/auth"
	"github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

	dummyDoer := testsupport.NewDummyHttpDoer()
	dummySpaceCollabCollector := &testsupport.DummyCollaboratorCollector{}
	dummyCollabResult, _ := dummySpaceCollabCollector.GetSpaceCollaborators(context.Background(), nil, uuid.NewV4())
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(dummyCollabResult)
	responseBody := ioutil.NopCloser(bytes.NewReader(b.Bytes()))

	dummyDoer.Client.Response = &http.Response{Body: responseBody, StatusCode: http.StatusOK}

	c := api.New(dummyDoer)
	id, _ := uuid.FromString("020f756e-b51a-4b43-b113-45cec16b9ce9")

	collector := auth.AuthCollector{}

	u, err := collector.GetSpaceCollaborators(context.Background(), c, id)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, u.Data, 10)

	// there is no easy way to test if the sa token really gets to override the privacy of emails.
	// the following lines only checks whether the emails show up at all if privacy is set to true
	for _, user := range u.Data {
		require.NotNil(t, user.Attributes.Email)
		if user.Attributes.EmailPrivate != nil && *user.Attributes.EmailPrivate {
			assert.Empty(t, *user.Attributes.Email)
		} else {
			assert.NotEmpty(t, *user.Attributes.Email)
		}
	}

}
