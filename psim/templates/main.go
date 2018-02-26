package templates

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"html/template"
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
