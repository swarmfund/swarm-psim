package request_monitor

import "time"

type Config struct {
	RequestTimeout time.Duration `fig:"request_timeout,required"`
	SleepPeriod    time.Duration `fig:"sleep_period"`
}
