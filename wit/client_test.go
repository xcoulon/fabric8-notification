package wit_test

import (
	"context"
	"testing"

	"github.com/fabric8-services/fabric8-notification/wit"
	"github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	OpenshiftIOAPI = "http://api.openshift.io"
)

func createClient(t *testing.T) *api.Client {
	c, err := wit.NewCachedClient(OpenshiftIOAPI)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestGetUser(t *testing.T) {

	c := createClient(t)
	ID, err := uuid.FromString("b67f1cee-0a9f-40da-8e52-504c092e54e0")

	u, err := wit.GetUser(context.Background(), c, ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "aslak@redhat.com", *u.Data.Attributes.Email)
}

func TestWorkItem(t *testing.T) {

	c := createClient(t)
	ID, err := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	u, err := wit.GetWorkItem(context.Background(), c, ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Cannot resolve Area/Iteration info for new WI created in in-memory mode", u.Data.Attributes["system.title"])
}

func TestComment(t *testing.T) {

	c := createClient(t)
	ID, err := uuid.FromString("5e7c1da9-af62-4b73-b18a-e88b7a6b9054")

	u, err := wit.GetComment(context.Background(), c, ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "The fields for Area/Iteration just shows loading animation, when the new WI is opened.", *u.Data.Attributes.Body)
}

func TestComments(t *testing.T) {

	c := createClient(t)
	id, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	u, err := wit.GetComments(context.Background(), c, id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "The fields for Area/Iteration just shows loading animation, when the new WI is opened.", *u.Data[0].Attributes.Body)
}

func TestSpace(t *testing.T) {

	c := createClient(t)
	id, _ := uuid.FromString("020f756e-b51a-4b43-b113-45cec16b9ce9")

	u, err := wit.GetSpace(context.Background(), c, id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "openshiftio", *u.Data.Attributes.Name)
}

func TestSpaceCollaborators(t *testing.T) {

	c := createClient(t)
	id, _ := uuid.FromString("020f756e-b51a-4b43-b113-45cec16b9ce9")

	u, err := wit.GetSpaceCollaborators(context.Background(), c, id)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, len(u.Data) > 10)
}
