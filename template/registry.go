package template

import (
	"bytes"
	"fmt"
	tmplHTML "html/template"
	tmplTXT "text/template"
)

type Template struct {
	Name    string
	Subject string
	Body    string
}

func (t Template) Render(vars map[string]interface{}) (subject string, body string, err error) {
	body, err = t.renderHTML(t.Body, vars)
	if err != nil {
		return "", "", err
	}
	subject, err = t.renderTXT(t.Subject, vars)
	if err != nil {
		return "", "", err
	}
	return subject, body, nil
}

func (t Template) renderHTML(template string, vars map[string]interface{}) (string, error) {
	funcMap := tmplHTML.FuncMap{
		"date":      formatDate,
		"sizeImage": sizeImage,
		"raw": func(s interface{}) tmplHTML.HTML {
			if s == nil {
				return tmplHTML.HTML("")
			}
			return tmplHTML.HTML(resolveString(s))
		},
		"lower": lower,
	}

	templ, err := tmplHTML.New(t.Name).Funcs(funcMap).Parse(template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = templ.Execute(&buf, vars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t Template) renderTXT(template string, vars map[string]interface{}) (string, error) {
	funcMap := tmplTXT.FuncMap{
		"date":  formatDate,
		"lower": lower,
	}

	templ, err := tmplTXT.New(t.Name).Funcs(funcMap).Parse(template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = templ.Execute(&buf, vars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type Registry interface {
	Get(string) (Template, bool)
}

type AssetRegistry struct {
}

func (a *AssetRegistry) Get(name string) (Template, bool) {
	bpath := fmt.Sprintf("template/%s/body.html", name)
	b, err := Asset(bpath)
	if err != nil {
		return Template{}, false
	}
	spath := fmt.Sprintf("template/%s/subject.txt", name)
	s, err := Asset(spath)
	if err != nil {
		return Template{}, false
	}

	return Template{
		Name:    name,
		Subject: string(s),
		Body:    string(b),
	}, true
}
