package collector_test

import (
	"context"
	"testing"

	"github.com/fabric8-services/fabric8-notification/auth"
	authApi "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/wit"
	witApi "github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	OpenshiftIOAPI     = "http://api.openshift.io"
	OpenShiftIOAuthAPI = "https://auth.openshift.io"
)

func createClient(t *testing.T) (*witApi.Client, *authApi.Client) {
	c, err := wit.NewCachedClient(OpenshiftIOAPI)
	if err != nil {
		t.Fatal(err)
	}

	authApi, err := auth.NewCachedClient(OpenShiftIOAuthAPI)
	if err != nil {
		t.Fatal(err)
	}
	return c, authApi
}

func TestWorkItem(t *testing.T) {

	witClient, authClient := createClient(t)
	wiID, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	users, vars, err := collector.WorkItem(context.Background(), authClient, witClient, wiID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vars["workitem"])
	assert.NotNil(t, vars["workitemOwner"])
	assert.NotNil(t, vars["workitemArea"])
	assert.NotNil(t, vars["workitemType"])
	assert.NotNil(t, vars["space"])
	assert.NotNil(t, vars["spaceOwner"])

	assert.True(t, len(users) > 10)

	/*
		for _, u := range users {
			fmt.Println(u.FullName, u.EMail)
		}
	*/
}

func TestComment(t *testing.T) {

	witClient, authClient := createClient(t)
	cID, _ := uuid.FromString("5e7c1da9-af62-4b73-b18a-e88b7a6b9054")

	users, vars, err := collector.Comment(context.Background(), authClient, witClient, cID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vars["comment"])
	assert.NotNil(t, vars["commentOwner"])
	assert.NotNil(t, vars["workitem"])
	assert.NotNil(t, vars["workitemOwner"])
	assert.NotNil(t, vars["workitemArea"])
	assert.NotNil(t, vars["workitemType"])
	assert.NotNil(t, vars["space"])
	assert.NotNil(t, vars["spaceOwner"])

	assert.True(t, len(users) > 10)

	/*
		for _, u := range users {
			fmt.Println(u.FullName, u.EMail)
		}
	*/
}
