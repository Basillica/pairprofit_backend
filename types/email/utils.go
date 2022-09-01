package email

import (
	"io/ioutil"
	"strings"
)

type EmailTemplate struct {
	Html              string
	TemplateFormatMap map[string]string
}

func (e *EmailTemplate) ParseHtmlFileToString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (e *EmailTemplate) FormatTemplateString() string {
	for k, v := range e.TemplateFormatMap {
		e.Html = strings.Replace(e.Html, "{{ "+k+" }}", v, 1)
	}
	return e.Html
}
