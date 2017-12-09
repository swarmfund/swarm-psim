package emails

import (
	"fmt"
	"html/template"
	"path"
	"strings"

	"gitlab.com/distributed_lab/logan/v3"
)

//go:generate go-bindata -ignore .+\.go$ -pkg emails -o bindata.go ./...
//go:generate gofmt -s -w bindata.go

const (
	templatesDir = "templates"
)

type AssetFn func(name string) ([]byte, error)

type AssetDirFn func(name string) ([]string, error)

type TemplatesLoader struct {
	asset    AssetFn
	assetDir AssetDirFn
	template *template.Template
}

var (
	Templates *TemplatesLoader
)

func init() {
	Templates = NewTemplatesLoader()
	if err := Templates.loadDir(templatesDir); err != nil {
		logan.New().
			WithField("service", "load-templates").
			WithError(err).
			Fatal("failed to load templates")
		return
	}
}

func NewTemplatesLoader() *TemplatesLoader {
	return &TemplatesLoader{
		asset:    Asset,
		assetDir: AssetDir,
		template: template.New("templates"),
	}
}

func (t *TemplatesLoader) loadDir(dir string) error {
	files, err := t.assetDir(dir)
	if err != nil {
		return err
	}

	for _, fp := range files {
		looksLikeTemplate := strings.HasSuffix(fp, ".html")
		if !looksLikeTemplate {
			t.loadDir(path.Join(dir, fp))
			continue
		}
		name := path.Join(dir, fp)
		bytes, err := t.asset(name)
		if err != nil {
			return err
		}
		_, err = t.template.New(name).Parse(string(bytes))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TemplatesLoader) Lookup(name string) *template.Template {
	name = fmt.Sprintf("%s/%s.html", templatesDir, name)
	return t.template.Lookup(name)
}
