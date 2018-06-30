package request_monitor

import "time"

type Config struct {
	RequestTimeout time.Duration `fig:"request_timeout,required"`
	AbnormalPeriodMin time.Duration `fig:"abnormal_period_min"`
	AbnormalPeriodMax time.Duration `fig:"abnormal_period_max"`
	SleepPeriod time.Duration `fig:"sleep_period"`
}
