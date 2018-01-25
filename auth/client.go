package auth

import (
	"net/http"
	"net/url"

	"github.com/fabric8-services/fabric8-notification/auth/api"
	goaclient "github.com/goadesign/goa/client"
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
