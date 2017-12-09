package q

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"gitlab.com/distributed_lab/notificator-server/log"
	_ "gopkg.in/mattes/migrate.v1/driver/postgres" // driver for migrations
	"gopkg.in/mattes/migrate.v1/migrate"
)

var instance Interface

type Interface interface {
	DB() *sqlx.DB
	Request() RequestQInterface
	Auth() AuthQInterface
}

type Q struct {
	db *sqlx.DB
}

func NewQ() Interface {
	entry := log.WithField("service", "q")
	db, err := sqlx.Open(conf.Driver, conf.DSN)
	if err != nil {
		entry.WithError(err).Fatal()
	}

	err = db.Ping()
	if err != nil {
		entry.WithError(err).Fatal()
	}

	return &Q{
		db: db,
	}
}

func (q *Q) Request() RequestQInterface {
	return &RequestQ{
		db: instance.DB(),
	}
}

func (q *Q) Auth() AuthQInterface {
	return &AuthQ{
		db: instance.DB(),
	}
}

func (q *Q) DB() *sqlx.DB {
	return q.db
}

func Init() {
	// highly synchronous procedure
	instance = NewQ()
}

func Request() RequestQInterface {
	return instance.Request()
}

func Migrate(migrations string) {
	entry := log.WithField("service", "migrate")
	errs, ok := migrate.UpSync(conf.DSN, migrations)
	if !ok {
		for _, err := range errs {
			entry.WithError(err).Error()
		}
		entry.Fatal("failed to migrate")
	}
	entry.Info("migrated successfully")
}

func NewMigration(migrations, name string) error {
	_, err := migrate.Create(conf.DSN, migrations, name)
	return err
}
