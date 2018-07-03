package testsupport

import (
	"context"
	authApi "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/goadesign/goa/uuid"
)

type DummyCollaboratorCollector struct {
}

func (c *DummyCollaboratorCollector) GetSpaceCollaborators(ctx context.Context, client *authApi.Client, spaceID uuid.UUID) (*authApi.UserList, error) {
	users := authApi.UserList{
		Data: []*authApi.UserData{},
	}
	for i := 0; i < 5; i++ {
		users.Data = append(users.Data, generateUserData(false))
		users.Data = append(users.Data, generateUserData(true))
	}

	return &users, nil
}

func generateUserData(verifiedEmail bool) *authApi.UserData {
	id := uuid.NewV4().String()
	var email string
	if verifiedEmail {
		email = "verified" + id
	} else {
		email = "unverified" + id
	}
	data := authApi.UserData{
		ID:   &id,
		Type: "identities",
		Attributes: &authApi.UserDataAttributes{
			Username:      &id,
			Email:         &email,
			FullName:      &id,
			EmailVerified: &verifiedEmail,
		},
	}
	return &data
}
