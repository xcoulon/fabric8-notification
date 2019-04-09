package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fabric8-services/fabric8-notification/auth"
	authapi "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/fabric8-services/fabric8-notification/testsupport"
	"github.com/fabric8-services/fabric8-notification/types"
	"github.com/fabric8-services/fabric8-notification/wit"
	"github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/goadesign/goa/uuid"
)

const (
	OpenshiftIOAPI     = "https://api.openshift.io"
	AuthOpenShiftIOAPI = "https://auth.openshift.io"
)

func main() {
	c, err := wit.NewCachedClient(OpenshiftIOAPI)
	authClient, err := auth.NewCachedClient(AuthOpenShiftIOAPI)
	if err != nil {
		panic(err)
	}

	type data struct {
		id           string
		templateName string
	}

	testdata := []data{
		{"de4871ce-0bfd-4b4b-aee2-e02427f4e38b", string(types.WorkitemCreate)},
		{"43024450-fe8c-4082-8828-88512cebfdb0", string(types.WorkitemCreate)},
		{"3a331aa3-6423-4fd7-85e4-95d7932b168c", string(types.WorkitemCreate)},
		{"d85e19a1-f4aa-486e-a8fe-3211cac9b68f", string(types.WorkitemCreate)},
		{"43024450-fe8c-4082-8828-88512cebfdb0", string(types.WorkitemUpdate)},
		{"d28f8344-4956-497a-b43b-7f217087a931", string(types.CommentCreate)},
		{"51d968b1-b9e5-4ec1-884a-ff256902c753", string(types.CommentCreate)},
		{"51d968b1-b9e5-4ec1-884a-ff256902c753", string(types.CommentUpdate)},
		{"3383826c-51e4-401b-9ccd-b898f7e2397d", string(types.UserEmailUpdate)},
		{"f3a184bf-963b-448a-86b4-307ae455223b", string(types.UserDeactivation)},
		{"81d1c3bf-fcf2-4c4e-9d12-f9e5c15fb9ab", string(types.InvitationTeamNoorg)},
		{"297f2037-72e9-42b3-a5fc-76d843877163", string(types.InvitationSpaceNoorg)},
		{"0a9c6814-462e-411c-8560-d74297bf1ceb", string(types.AnalyticsNotifyCVE)},
		{"7dd737ed-c0fa-4f2f-9a3a-ed43778bb7fe", string(types.AnalyticsNotifyVersion)},
	}
	fmt.Println("Generating test templates..")
	fmt.Println("")

	for _, d := range testdata {
		err = generate(authClient, c, d.id, d.templateName)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}

func generate(authClient *authapi.Client, c *api.Client, id, tmplName string) error {
	reg := template.AssetRegistry{}

	tmpl, exist := reg.Get(tmplName)
	if !exist {
		return fmt.Errorf("template %v not found", tmplName)
	}

	wiID, _ := uuid.FromString(id)

	var vars map[string]interface{}
	var err error

	// When running locally the actor ID has to be mocked
	// since there is no real actor.
	OSIOctx := testsupport.CreateOSIOUserContext()
	if strings.HasPrefix(tmplName, "workitem") {
		_, vars, err = collector.WorkItem(OSIOctx, authClient, c, nil, wiID)
	} else if strings.HasPrefix(tmplName, "comment") {
		_, vars, err = collector.Comment(OSIOctx, authClient, c, nil, wiID)
	} else if strings.HasPrefix(tmplName, "user.deactivation") {
		vars["custom"] = map[string]interface{}{
			"userEmail":  "user@example.com",
			"expiryDate": time.Now().Add(7 * 24 * time.Hour).Format("Jan 2, 2006"),
		}
	} else if strings.HasPrefix(tmplName, "user") {
		_, vars, err = collector.User(context.Background(), authClient, wiID)
		vars["custom"] = map[string]interface{}{
			// a realistic verifyURL
			"verifyURL": "https://auth.prod-preview.openshift.io/api/users/verifyemail?code=580f7d71-853c-48df-8206-d1265bcf44f1",
		}
	} else if strings.HasPrefix(tmplName, "invitation") {
		vars = make(map[string]interface{})
		vars["custom"] = map[string]interface{}{
			"inviter":   "Albert Einstein",
			"spaceName": "Physics Research Club",
			"teamName":  "Temporal Dynamics",
			"roleNames": "Scientist, Researcher",
			"acceptURL": "http://openshift.io/invitations/accept/12345-ABCDE-FFFFF-99999-88888",
		}

	} else if strings.HasPrefix(tmplName, "analytics.notify.cve") {
		vars = make(map[string]interface{})
		payload, err := testsupport.GetFileContent("preview/test-files/cve.payload.json")
		if err == nil {
			vars["custom"] = testsupport.GetCustomElement(payload)
		}
	} else if strings.HasPrefix(tmplName, "analytics.notify.version") {
		vars = make(map[string]interface{})
		payload, err := testsupport.GetFileContent("preview/test-files/version.payload.json")
		if err == nil {
			vars["custom"] = testsupport.GetCustomElement(payload)
		}
	} else {
		return fmt.Errorf("Unkown resolver for template %v", tmplName)
	}

	if err != nil {
		if len(vars) == 0 {
			return err
		}
	}

	fileName, err := filepath.Abs("tmp/" + tmplName + "-" + id + ".html")
	if err != nil {
		return err
	}
	subject, body, headers, err := tmpl.Render(addGlobalVars(vars))
	if err != nil {
		return err
	}
	fmt.Println("Subject:", subject)
	fmt.Println("Output :", "file://"+fileName)
	fmt.Println("Headers:")
	for k, v := range headers {
		fmt.Println(k, v)
	}
	fmt.Println("")

	ioutil.WriteFile(fileName, []byte(body), os.FileMode(0777))
	return nil
}

func addGlobalVars(vars map[string]interface{}) map[string]interface{} {
	vars["webURL"] = "https://openshift.io"
	return vars
}
