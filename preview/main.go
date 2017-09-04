package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/fabric8-services/fabric8-notification/wit"
	"github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/goadesign/goa/uuid"
)

const (
	OpenshiftIOAPI = "https://api.openshift.io"
)

func main() {
	c, err := wit.NewCachedClient(OpenshiftIOAPI)
	if err != nil {
		panic(err)
	}

	type data struct {
		id           string
		templateName string
	}

	var testdata []data
	testdata = append(testdata, data{"43024450-fe8c-4082-8828-88512cebfdb0", "workitem.create"})
	testdata = append(testdata, data{"3a331aa3-6423-4fd7-85e4-95d7932b168c", "workitem.create"})
	testdata = append(testdata, data{"d85e19a1-f4aa-486e-a8fe-3211cac9b68f", "workitem.create"})
	testdata = append(testdata, data{"43024450-fe8c-4082-8828-88512cebfdb0", "workitem.update"})

	testdata = append(testdata, data{"d28f8344-4956-497a-b43b-7f217087a931", "comment.create"})
	testdata = append(testdata, data{"51d968b1-b9e5-4ec1-884a-ff256902c753", "comment.create"})
	testdata = append(testdata, data{"51d968b1-b9e5-4ec1-884a-ff256902c753", "comment.update"})

	fmt.Println("Generating test templates..")
	fmt.Println("")

	for _, d := range testdata {
		err = generate(c, d.id, d.templateName)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}

func generate(c *api.Client, id, tmplName string) error {
	reg := template.AssetRegistry{}

	temp, exist := reg.Get(tmplName)
	if !exist {
		return fmt.Errorf("Tempalte %v not found", tmplName)
	}

	wiID, _ := uuid.FromString(id)

	var vars map[string]interface{}
	var err error

	if strings.HasPrefix(tmplName, "workitem") {
		_, vars, err = collector.WorkItem(context.Background(), c, wiID)
	} else if strings.HasPrefix(tmplName, "comment") {
		_, vars, err = collector.Comment(context.Background(), c, wiID)
	} else {
		return fmt.Errorf("Unkown resovler for template %v", tmplName)
	}

	if err != nil {
		return err
	}

	fileName, err := filepath.Abs("tmp/" + tmplName + "-" + id + ".html")
	if err != nil {
		return err
	}
	subject, body, _, err := temp.Render(addGlobalVars(vars))
	if err != nil {
		return err
	}
	fmt.Println("Subject:", subject)
	fmt.Println("Output :", "file://"+fileName)

	ioutil.WriteFile(fileName, []byte(body), os.FileMode(0777))
	return nil
}

func addGlobalVars(vars map[string]interface{}) map[string]interface{} {
	vars["webURL"] = "https://openshift.io"
	return vars
}
