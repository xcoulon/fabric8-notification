package template_test

import (
	"context"
	"testing"

	"strings"

	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/fabric8-services/fabric8-notification/wit"
	"github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	OpenshiftIOAPI = "http://api.openshift.io"
)

func addGlobalVars(vars map[string]interface{}) map[string]interface{} {
	vars["webURL"] = "http://localhost"
	return vars
}

func createClient(t *testing.T) *api.Client {
	c, err := wit.NewCachedClient(OpenshiftIOAPI)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestTrueOnFoundName(t *testing.T) {
	reg := template.AssetRegistry{}

	_, exist := reg.Get("workitem.update")
	assert.True(t, exist)
}

func TestFalseOnMissingName(t *testing.T) {
	reg := template.AssetRegistry{}

	_, exist := reg.Get("MISSING")
	assert.False(t, exist)
}

func TestRenderWorkitemCreate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get("workitem.create")
	assert.True(t, exist)

	c := createClient(t)
	wiID, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	_, vars, err := collector.WorkItem(context.Background(), c, wiID)
	if err != nil {
		t.Fatal(err)
	}

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)
	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "Message-ID")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "Ruchir Garg"))
	assert.True(t, strings.Contains(body, "1343"))
	assert.True(t, strings.Contains(body, "mode under Backlog")) // part of the Description. Might change since we're on live data

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderWorkitemCreateMissingDescription(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get("workitem.create")
	assert.True(t, exist)

	c := createClient(t)
	wiID, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	_, vars, err := collector.WorkItem(context.Background(), c, wiID)
	if err != nil {
		t.Fatal(err)
	}

	wi := vars["workitem"].(*api.WorkItemSingle)
	delete(wi.Data.Attributes, "system.description")
	delete(wi.Data.Attributes, "system.description.rendered")
	delete(wi.Data.Attributes, "system.description.markup")

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)
	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "Message-ID")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "Ruchir Garg"))
	assert.True(t, strings.Contains(body, "1343"))

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderWorkitemUpdate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get("workitem.update")
	assert.True(t, exist)

	c := createClient(t)
	wiID, _ := uuid.FromString("8bccc228-bba7-43ad-b077-15fbb9148f7f")

	_, vars, err := collector.WorkItem(context.Background(), c, wiID)
	if err != nil {
		t.Fatal(err)
	}

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)

	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "In-Reply-To")
	assert.Contains(t, headers, "References")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "1343"))

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderCommentCreate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get("comment.create")
	assert.True(t, exist)

	c := createClient(t)
	ciID, _ := uuid.FromString("5e7c1da9-af62-4b73-b18a-e88b7a6b9054")

	_, vars, err := collector.Comment(context.Background(), c, ciID)
	if err != nil {
		t.Fatal(err)
	}

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)

	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "In-Reply-To")
	assert.Contains(t, headers, "References")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "1343"))
	assert.True(t, strings.Contains(body, "just shows loading animation,")) // part of the msg

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}

func TestRenderCommentUpdate(t *testing.T) {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get("comment.update")
	assert.True(t, exist)

	c := createClient(t)
	ciID, _ := uuid.FromString("5e7c1da9-af62-4b73-b18a-e88b7a6b9054")

	_, vars, err := collector.Comment(context.Background(), c, ciID)
	if err != nil {
		t.Fatal(err)
	}

	subject, body, headers, err := temp.Render(addGlobalVars(vars))
	require.NoError(t, err)

	assert.True(t, strings.Contains(subject, "[openshiftio/openshiftio]"))
	assert.True(t, strings.Contains(subject, "[Scenario]"))

	assert.Contains(t, headers, "In-Reply-To")
	assert.Contains(t, headers, "References")

	assert.True(t, strings.Contains(body, "http://localhost/openshiftio/openshiftio/plan/detail/1343"))
	assert.True(t, strings.Contains(body, "1343"))
	assert.True(t, strings.Contains(body, "just shows loading animation,")) // part of the ms

	/*
		ioutil.WriteFile("../test.html", []byte(body), os.FileMode(0777))
	*/
}
