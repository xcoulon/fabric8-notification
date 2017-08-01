package template

import (
	"bytes"
	"fmt"
	tmplHTML "html/template"
	"reflect"
	"strings"
	tmplTXT "text/template"
	"time"
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
		"date": formatDate,
		"raw": func(s interface{}) tmplHTML.HTML {
			if s == nil {
				return tmplHTML.HTML("")
			}
			return tmplHTML.HTML(fmt.Sprint(s))
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

func formatDate(date interface{}) string {
	format := "02 January 2006"
	if d, ok := date.(*time.Time); ok {
		return d.Format(format)
	}
	if d, ok := date.(string); ok {
		p, err := time.Parse(time.RFC3339, d)
		if err != nil {
			return "Unknown"
		}
		return p.Format(format)
	}
	return "Unknown"
}

func lower(data interface{}) string {
	if data == nil {
		return ""
	}

	var d string

	value := reflect.ValueOf(data)
	if value.Type().Kind() == reflect.Ptr {
		d = fmt.Sprint(value.Elem())
	} else {
		d = fmt.Sprint(data)
	}
	return strings.ToLower(d)
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
