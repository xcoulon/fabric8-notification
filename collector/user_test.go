package collector_test

import (
	"context"
	"testing"

	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	c := createClient(t)
	uID, _ := uuid.FromString("3383826c-51e4-401b-9ccd-b898f7e2397d")
	users, vars, err := collector.User(context.Background(), c, uID)

	assert.Nil(t, err)
	assert.Len(t, users, 1)
	assert.Len(t, vars, 0)
}
