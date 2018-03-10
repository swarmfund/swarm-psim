package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type DBConf struct {
	Driver string
	DSN    string
}

const dbConfigKey = "db"

func (c *ViperConfig) DB() DBConf {
	c.Lock()
	defer c.Unlock()

	if c.db != nil {
		return *c.db
	}

	db := new(DBConf)
	config := c.GetStringMap(dbConfigKey)

	if err := figure.Out(db).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out db"))
	}

	c.db = db

	return *c.db
}
