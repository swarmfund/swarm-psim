package conf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type S3 struct {
	Url              string `fig:"url,required"`
	AccessKey        string `fig:"access_key,required"`
	SecretKey        string `fig:"secret_key,required"`
	Region           string `fig:"region"`
	DisableSSL       bool   `fig:"disable_ssl"`
	S3ForcePathStyle bool   `fig:"s3_force_path_style"`
}

func (c *ViperConfig) S3() *session.Session {
	c.Lock()
	defer c.Unlock()

	if c.session != nil {
		return c.session
	}

	provider := &S3{}
	config := c.Get(ServiceS3)

	if err := figure.Out(provider).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out S3"))
	}

	sess, err := CreateSession(*provider)
	if err != nil {
		panic(errors.Wrap(err, "failed to create session S3"))
	}

	c.session = sess

	return c.session
}

func CreateSession(provider S3) (*session.Session, error) {
	cfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(provider.AccessKey, provider.SecretKey, ""),
		Endpoint:         aws.String(provider.Url),
		Region:           aws.String(provider.Region),
		DisableSSL:       aws.Bool(provider.DisableSSL),
		S3ForcePathStyle: aws.Bool(provider.S3ForcePathStyle),
	}

	return session.NewSession(cfg)
}
