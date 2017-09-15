package collector

import (
	"context"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fabric8-services/fabric8-notification/configuration"
	"github.com/fabric8-services/fabric8-notification/wit"
	"github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/fabric8-services/fabric8-wit/log"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/goadesign/goa/uuid"
)

func NewCommentResolver(c *api.Client) ReceiverResolver {
	return func(ctx context.Context, id string) ([]Receiver, map[string]interface{}, error) {
		cID, err := uuid.FromString(id)
		if err != nil {
			return []Receiver{}, nil, fmt.Errorf("unable to lookup comment based on id %v", id)
		}
		return Comment(ctx, c, cID)
	}
}

func NewWorkItemResolver(c *api.Client) ReceiverResolver {
	return func(ctx context.Context, id string) ([]Receiver, map[string]interface{}, error) {
		wID, err := uuid.FromString(id)
		if err != nil {
			return []Receiver{}, nil, fmt.Errorf("unable to lookup Workitem based on id %v", id)
		}
		return WorkItem(ctx, c, wID)
	}
}

func ConfiguredVars(config *configuration.Data, resolver ReceiverResolver) ReceiverResolver {
	return func(ctx context.Context, id string) ([]Receiver, map[string]interface{}, error) {
		r, v, err := resolver(ctx, id)
		if err != nil {
			return r, v, err
		}

		v["webURL"] = config.GetWebURL()
		return r, v, err
	}
}

func Comment(ctx context.Context, c *api.Client, cID uuid.UUID) ([]Receiver, map[string]interface{}, error) {
	var values = map[string]interface{}{}
	var errors []error
	var users []uuid.UUID

	comment, err := wit.GetComment(ctx, c, cID)
	if err != nil {
		return []Receiver{}, nil, err
	}
	values["comment"] = comment
	users = append(users, collectCommentUsers(comment)...)

	commentOwner, err := wit.GetUser(ctx, c, *comment.Data.Relationships.CreatedBy.Data.ID)
	if err != nil {
		errors = append(errors, err)
	}
	values["commentOwner"] = commentOwner

	wiID, _ := uuid.FromString(*comment.Data.Relationships.Parent.Data.ID)

	wi, err := wit.GetWorkItem(ctx, c, wiID)
	if err != nil {
		return []Receiver{}, nil, err
	}
	users = append(users, collectWorkItemUsers(wi)...)
	values["workitem"] = wi

	ownerID, _ := uuid.FromString(*wi.Data.Relationships.Creator.Data.ID)
	workitemOwner, err := wit.GetUser(ctx, c, ownerID)
	if err != nil {
		errors = append(errors, err)
	}
	values["workitemOwner"] = workitemOwner

	areaID, _ := uuid.FromString(*wi.Data.Relationships.Area.Data.ID)
	workitemArea, err := wit.GetArea(ctx, c, areaID)
	if err != nil {
		errors = append(errors, err)
	}
	values["workitemArea"] = workitemArea

	cs, err := wit.GetComments(ctx, c, wiID)
	if err != nil {
		errors = append(errors, err)
	}
	users = append(users, collectCommentsUsers(cs)...)

	spaceID := *wi.Data.Relationships.Space.Data.ID
	s, err := wit.GetSpace(ctx, c, spaceID)
	if err != nil {
		errors = append(errors, err)
	}
	users = append(users, collectSpaceUsers(s)...)
	values["space"] = s

	spaceOwner, err := wit.GetUser(ctx, c, *s.Data.Relationships.OwnedBy.Data.ID)
	if err != nil {
		errors = append(errors, err)
	}
	values["spaceOwner"] = spaceOwner

	witype, err := wit.GetWorkItemType(ctx, c, spaceID, wi.Data.Relationships.BaseType.Data.ID)
	if err != nil {
		errors = append(errors, err)
	}
	values["workitemType"] = witype

	actorID, err := getActorID(ctx)
	if err == nil {
		actor, err := wit.GetUser(ctx, c, actorID)
		if err != nil {
			errors = append(errors, err)
		}
		values["actor"] = actor
	}

	sc, err := wit.GetSpaceCollaborators(ctx, c, spaceID)
	if err != nil {
		errors = append(errors, err)
	}
	users = append(users, collectSpaceCollaboratorUsers(sc)...)

	resolved, err := resolveAllUsers(ctx, c, SliceUniq(users), sc.Data)
	if err != nil {
		errors = append(errors, err)
	}
	resolved = removeActorFromReceivers(ctx, resolved)

	if len(errors) > 0 {
		return resolved, values, multiError{Message: "errors during notification resolving", Errors: errors}
	}

	return resolved, values, nil
}

func WorkItem(ctx context.Context, c *api.Client, wiID uuid.UUID) ([]Receiver, map[string]interface{}, error) {
	var values = map[string]interface{}{}
	var errors []error
	var users []uuid.UUID

	wi, err := wit.GetWorkItem(ctx, c, wiID)
	if err != nil {
		return []Receiver{}, nil, err
	}
	values["workitem"] = wi
	users = append(users, collectWorkItemUsers(wi)...)

	ownerID, _ := uuid.FromString(*wi.Data.Relationships.Creator.Data.ID)
	workitemOwner, err := wit.GetUser(ctx, c, ownerID)
	if err != nil {
		errors = append(errors, err)
	}
	values["workitemOwner"] = workitemOwner

	areaID, _ := uuid.FromString(*wi.Data.Relationships.Area.Data.ID)
	workitemArea, err := wit.GetArea(ctx, c, areaID)
	if err != nil {
		errors = append(errors, err)
	}
	values["workitemArea"] = workitemArea

	cs, err := wit.GetComments(ctx, c, wiID)
	if err != nil {
		errors = append(errors, err)
	}
	users = append(users, collectCommentsUsers(cs)...)

	spaceID := *wi.Data.Relationships.Space.Data.ID
	s, err := wit.GetSpace(ctx, c, spaceID)
	if err != nil {
		errors = append(errors, err)
	}
	values["space"] = s
	users = append(users, collectSpaceUsers(s)...)

	spaceOwner, err := wit.GetUser(ctx, c, *s.Data.Relationships.OwnedBy.Data.ID)
	if err != nil {
		errors = append(errors, err)
	}
	values["spaceOwner"] = spaceOwner

	witype, err := wit.GetWorkItemType(ctx, c, spaceID, wi.Data.Relationships.BaseType.Data.ID)
	if err != nil {
		errors = append(errors, err)
	}
	values["workitemType"] = witype

	actorID, err := getActorID(ctx)
	if err == nil {
		actor, err := wit.GetUser(ctx, c, actorID)
		if err != nil {
			errors = append(errors, err)
		}
		values["actor"] = actor
	}

	sc, err := wit.GetSpaceCollaborators(ctx, c, spaceID)
	if err != nil {
		errors = append(errors, err)
	}
	users = append(users, collectSpaceCollaboratorUsers(sc)...)

	resolved, err := resolveAllUsers(ctx, c, SliceUniq(users), sc.Data)
	if err != nil {
		errors = append(errors, err)
	}
	resolved = removeActorFromReceivers(ctx, resolved)

	if len(errors) > 0 {
		return resolved, values, multiError{Message: "errors during notification resolving", Errors: errors}
	}

	return resolved, values, nil
}

func resolveAllUsers(ctx context.Context, c *api.Client, users []uuid.UUID, collaborators []*api.UserData) ([]Receiver, error) {
	var resolved []Receiver

	for _, u := range users {
		found := false
		for _, c := range collaborators {
			if u.String() == *c.ID {
				found = true
				if c.Attributes.Email != nil {
					user := Receiver{EMail: *c.Attributes.Email}
					if c.Attributes.FullName != nil {
						user.FullName = *c.Attributes.FullName
					}
					resolved = append(resolved, user)
				}
			}
		}
		if !found {
			usr, err := wit.GetUser(ctx, c, u)
			if err == nil {
				if usr.Data.Attributes.Email != nil {
					user := Receiver{EMail: *usr.Data.Attributes.Email}
					if usr.Data.Attributes.FullName != nil {
						user.FullName = *usr.Data.Attributes.FullName
					}
					resolved = append(resolved, user)
				}
			} else {
				log.Error(ctx, map[string]interface{}{
					"err":     err,
					"user_id": u,
				}, "unable to lookup user")

			}
		}
	}

	return resolved, nil
}

func collectSpaceCollaboratorUsers(cl *api.UserList) []uuid.UUID {
	var users []uuid.UUID
	for _, c := range cl.Data {
		cID, err := uuid.FromString(*c.ID)
		if err == nil {
			users = append(users, cID)
		}
	}
	return users
}

func collectCommentsUsers(cl *api.CommentList) []uuid.UUID {
	var users []uuid.UUID

	for _, c := range cl.Data {
		if c.Relationships.CreatedBy != nil {
			users = append(users, *c.Relationships.CreatedBy.Data.ID)
		}
	}

	return users
}

func collectCommentUsers(c *api.CommentSingle) []uuid.UUID {
	var users []uuid.UUID

	if c.Data.Relationships.CreatedBy != nil {
		users = append(users, *c.Data.Relationships.CreatedBy.Data.ID)
	}

	return users
}

func collectWorkItemUsers(wi *api.WorkItemSingle) []uuid.UUID {
	var users []uuid.UUID

	if wi.Data.Relationships.Creator != nil && wi.Data.Relationships.Creator.Data != nil {
		creatorID, err := uuid.FromString(*wi.Data.Relationships.Creator.Data.ID)
		if err == nil {
			users = append(users, creatorID)
		}
	}

	if wi.Data.Relationships.Assignees != nil && wi.Data.Relationships.Assignees.Data != nil {
		for _, assignee := range wi.Data.Relationships.Assignees.Data {
			assigneeID, err := uuid.FromString(*assignee.ID)
			if err == nil {
				users = append(users, assigneeID)
			}
		}
	}
	return users
}

func collectSpaceUsers(space *api.SpaceSingle) []uuid.UUID {
	var users []uuid.UUID

	if space.Data.Relationships.OwnedBy != nil {
		users = append(users, *space.Data.Relationships.OwnedBy.Data.ID)
	}
	return users
}

func removeActorFromReceivers(ctx context.Context, rec []Receiver) []Receiver {
	var res []Receiver
	actorEmail := getActorEmail(ctx)
	for _, rec := range rec {
		if rec.EMail != actorEmail {
			res = append(res, rec)
		}
	}
	return res
}

// SliceUniq removes duplicate values in given slice
func SliceUniq(a []uuid.UUID) []uuid.UUID {
	result := []uuid.UUID{}
	seen := map[uuid.UUID]uuid.UUID{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

// slightly missplaced.
func getActorID(ctx context.Context) (uuid.UUID, error) {
	token := goajwt.ContextJWT(ctx)
	if token == nil {
		return uuid.UUID{}, fmt.Errorf("Missing token")
	}
	id := token.Claims.(jwt.MapClaims)["sub"]
	if id == nil {
		return uuid.UUID{}, fmt.Errorf("Missing sub")
	}

	ID, err := uuid.FromString(id.(string))
	if err != nil {
		return uuid.UUID{}, err
	}
	return ID, nil
}

func getActorEmail(ctx context.Context) string {
	token := goajwt.ContextJWT(ctx)
	if token == nil {
		return ""
	}
	e := token.Claims.(jwt.MapClaims)["email"]
	if e == nil {
		return ""
	}

	email, err := e.(string)
	if !err {
		return ""
	}
	return email
}
