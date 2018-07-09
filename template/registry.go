package template

import (
	"bytes"
	"fmt"
	tmplHTML "html/template"
	tmplTXT "text/template"

	"github.com/magiconair/properties"
)

type Template struct {
	Name    string
	Subject string
	Body    string
	Headers string
}

func (t Template) Render(vars map[string]interface{}) (subject string, body string, headers map[string]string, err error) {
	body, err = t.renderHTML(t.Body, vars)
	if err != nil {
		return "", "", nil, err
	}
	subject, err = t.renderTXT(t.Subject, vars)
	if err != nil {
		return "", "", nil, err
	}
	headers, err = t.renderHeaders(t.Headers, vars)
	if err != nil {
		return "", "", nil, err
	}

	return subject, body, headers, err
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
		"lower":     lower,
		"detailURL": detailURL,
		"areaPath":  areaPath,
		"inc":       inc,
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
		"date":      formatDate,
		"lower":     lower,
		"detailURL": detailURL,
		"areaPath":  areaPath,
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

func (t Template) renderHeaders(template string, vars map[string]interface{}) (map[string]string, error) {
	h, err := t.renderTXT(template, vars)
	if err != nil {
		return nil, err
	}
	headers, err := properties.LoadString(h)
	if err != nil {
		return nil, err
	}
	return headers.Map(), nil
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
	hpath := fmt.Sprintf("template/%s/headers.prop", name)
	h, err := Asset(hpath)
	if err != nil {
		h = []byte{}
	}

	return Template{
		Name:    name,
		Subject: string(s),
		Body:    string(b),
		Headers: string(h),
	}, true
}
