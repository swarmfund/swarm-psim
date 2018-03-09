package q

import (
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"                          // postgres driver
	_ "gopkg.in/mattes/migrate.v1/driver/postgres" // driver for migrations
	"gopkg.in/mattes/migrate.v1/migrate"
)

var instance QInterface

type QInterface interface {
	DB() *sqlx.DB
	Request() RequestQInterface
	Auth() AuthQInterface
}

type Q struct {
	db *sqlx.DB
}

func NewQ(driver, dsn string, log *logrus.Logger) QInterface {
	entry := log.WithField("service", "q")

	db, err := sqlx.Open(driver, dsn)
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

func Init(driver, dsn string, log *logrus.Logger) {
	// highly synchronous procedure
	instance = NewQ(driver, dsn, log)
}

func GetQInstance() QInterface {
	return instance
}

func Request() RequestQInterface {
	return instance.Request()
}

func Migrate(dsn, migrations string, log *logrus.Logger) {
	entry := log.WithField("service", "migrate")
	errs, ok := migrate.UpSync(dsn, migrations)
	if !ok {
		for _, err := range errs {
			entry.WithError(err).Error()
		}
		entry.Fatal("failed to migrate")
	}
	entry.Info("migrated successfully")
}

func NewMigration(dsn, migrations, name string) error {
	_, err := migrate.Create(dsn, migrations, name)
	return err
}
