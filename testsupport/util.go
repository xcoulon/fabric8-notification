package testsupport

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

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
