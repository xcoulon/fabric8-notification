package testsupport

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	jwt "github.com/dgrijalva/jwt-go"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"

	"github.com/fabric8-services/fabric8-notification/app"
)

func GetFileContent(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetCustomElement(payload string) map[string]interface{} {
	notifyPayload := &app.SendNotifyPayload{}
	err := json.Unmarshal([]byte(payload), notifyPayload)
	if err != nil {
		return nil
	}
	return notifyPayload.Data.Attributes.Custom
}

// CreateOSIOUserContext creates a new context using "openshiftio" user's ID
func CreateOSIOUserContext() context.Context {
	claims := jwt.MapClaims{}
	// Set actor ID to "openshiftio" user ID.
	claims["sub"] = "7b50ddb4-5e12-4031-bca7-3b88f92e2339"
	ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))
	return ctx
}
