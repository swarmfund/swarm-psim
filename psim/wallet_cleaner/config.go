package wallet_cleaner

import (
	"time"
)

type Config struct {
	ExpireDuration time.Duration `fig:"expire_duration,required"`
}
