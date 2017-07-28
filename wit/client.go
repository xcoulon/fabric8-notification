package wit

import (
	"context"
	"net/http"
	"net/url"

	"fmt"

	"github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/fabric8-services/fabric8-wit/goasupport"
	goaclient "github.com/goadesign/goa/client"
	"github.com/goadesign/goa/uuid"
	"github.com/gregjones/httpcache"
)

func NewCachedClient(hostURL string) (*api.Client, error) {

	u, err := url.Parse(hostURL)
	if err != nil {
		return nil, err
	}

	tp := httpcache.NewMemoryCacheTransport()
	client := http.Client{Transport: tp}

	c := api.New(goaclient.HTTPClientDoer(&client))
	c.Host = u.Host
	c.Scheme = u.Scheme
	return c, nil
}

func GetUser(ctx context.Context, client *api.Client, uID uuid.UUID) (*api.User, error) {
	resp, err := client.ShowUsers(goasupport.ForwardContextRequestID(ctx), api.ShowUsersPath(uID.String()), nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET user", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeUser(resp)
}

func GetWorkItem(ctx context.Context, client *api.Client, wiID uuid.UUID) (*api.WorkItemSingle, error) {
	resp, err := client.ShowWorkitem(goasupport.ForwardContextRequestID(ctx), api.ShowWorkitemPath(wiID), nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET workitem", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeWorkItemSingle(resp)
}

func GetWorkItemType(ctx context.Context, client *api.Client, spaceID, witID uuid.UUID) (*api.WorkItemTypeSingle, error) {
	resp, err := client.ShowWorkitemtype(goasupport.ForwardContextRequestID(ctx), api.ShowWorkitemtypePath(spaceID, witID), nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET workitemtype", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeWorkItemTypeSingle(resp)
}

func GetComment(ctx context.Context, client *api.Client, cID uuid.UUID) (*api.CommentSingle, error) {
	resp, err := client.ShowComments(goasupport.ForwardContextRequestID(ctx), api.ShowCommentsPath(cID), nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET comment", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeCommentSingle(resp)
}

func GetComments(ctx context.Context, client *api.Client, wiID uuid.UUID) (*api.CommentList, error) {
	pageLimit := 100
	pageOffset := "0"

	resp, err := client.ListWorkItemComments(goasupport.ForwardContextRequestID(ctx), api.ListWorkItemCommentsPath(wiID), &pageLimit, &pageOffset, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET comments", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeCommentList(resp)
}

func GetSpace(ctx context.Context, client *api.Client, spaceID uuid.UUID) (*api.SpaceSingle, error) {
	resp, err := client.ShowSpace(goasupport.ForwardContextRequestID(ctx), api.ShowSpacePath(spaceID), nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET space", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeSpaceSingle(resp)
}

func GetSpaceCollaborators(ctx context.Context, client *api.Client, spaceID uuid.UUID) (*api.UserList, error) {
	pageLimit := 100
	pageOffset := "0"
	resp, err := client.ListCollaborators(goasupport.ForwardContextRequestID(ctx), api.ListCollaboratorsPath(spaceID), &pageLimit, &pageOffset, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET collaborators", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeUserList(resp)
}
