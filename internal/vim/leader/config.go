package leader

import (
	"time"

	"github.com/jorgerojas26/lazysql/models"
)

type Config struct {
	Timeout time.Duration
}

func DefaultConfig() Config {
	return Config{Timeout: DefaultTimeout}
}

func ConfigFromApp(appConfig *models.AppConfig) Config {
	if appConfig == nil || appConfig.LeaderTimeoutMs <= 0 {
		return DefaultConfig()
	}
	return Config{Timeout: time.Duration(appConfig.LeaderTimeoutMs) * time.Millisecond}
}
