package collector_test

import (
	"context"
	"testing"

	"github.com/fabric8-services/fabric8-notification/collector"
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

func TestWorkItem(t *testing.T) {

	c := createClient(t)
	wiID, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	users, vars, err := collector.WorkItem(context.Background(), c, wiID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vars["workitem"])
	assert.NotNil(t, vars["workitemOwner"])
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

	c := createClient(t)
	cID, _ := uuid.FromString("5e7c1da9-af62-4b73-b18a-e88b7a6b9054")

	users, vars, err := collector.Comment(context.Background(), c, cID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vars["comment"])
	assert.NotNil(t, vars["commentOwner"])
	assert.NotNil(t, vars["workitem"])
	assert.NotNil(t, vars["workitemOwner"])
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
