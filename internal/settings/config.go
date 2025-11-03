package settings

import "time"

type Config struct {
	Port int
	Env  string
	DB   struct {
		URL          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
}
