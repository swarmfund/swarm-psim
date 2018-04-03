package template_provider

type Config struct {
	Host               string `fig:"host"`
	Port               int    `fig:"port"`
	Bucket             string `fig:"bucket,required"`
	SkipSignatureCheck bool   `fig:"skip_signature_check"`
}
