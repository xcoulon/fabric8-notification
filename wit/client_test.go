package wit_test

import (
	"context"
	"testing"

	"github.com/fabric8-services/fabric8-notification/wit"
	"github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	openshiftIOAPI = "http://api.openshift.io"
)

func createClient(t *testing.T) *api.Client {
	c, err := wit.NewCachedClient(openshiftIOAPI)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestWorkItem(t *testing.T) {

	c := createClient(t)
	ID, err := uuid.FromString("4728edab-6ccb-4c99-bd4c-f4aeabc560ad")

	u, err := wit.GetWorkItem(context.Background(), c, ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "The default workitem type for inline quick-add in Experiences context is set to Bug instead of Feature", u.Data.Attributes["system.title"])
}

func TestArea(t *testing.T) {

	c := createClient(t)
	ID, err := uuid.FromString("b611ebaa-dfc9-489e-bb2f-c1d8d8237e40")

	u, err := wit.GetArea(context.Background(), c, ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Planner", *u.Data.Attributes.Name)
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
	require.NoError(t, err)

	require.NotNil(t, u)
	require.NotEmpty(t, u.Data)
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
