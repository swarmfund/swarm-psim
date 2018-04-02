package template_provider

const templateAPI = "template_api"

type TemplateAPI struct {
	Host               string `fig:"host"`
	Port               int    `fig:"port"`
	Bucket             string `fig:"bucket,required"`
	SkipSignatureCheck bool   `fig:"skip_signature_check"`
}
