package templates

import (
	"bytes"
	"html/template"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//go:generate go-bindata -nometadata -ignore .+\.go$ -pkg templates -o bindata.go ./...
//go:generate gofmt -w bindata.go

func GetHtmlTemplate(templateName string) (*template.Template, error) {
	bb, err := Asset(templateName)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to obtain template bytes")
	}

	t, err := template.New("template").Parse(string(bb))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse html.Template")
	}

	return t, nil
}

// TODO Cache the html template once and reuse it
func BuildTemplateEmailMessage(templateName string, templateData interface{}) (string, error) {
	fields := logan.F{
		"template_name": templateName,
	}

	t, err := GetHtmlTemplate(templateName)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get html Template", fields)
	}

	var buff bytes.Buffer

	err = t.Execute(&buff, templateData)
	if err != nil {
		return "", errors.Wrap(err, "Failed to execute html Template", fields)
	}

	return buff.String(), nil
}
