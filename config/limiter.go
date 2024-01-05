package config

import "golang.org/x/time/rate"

var Limiter = rate.NewLimiter(rate.Limit(10), 1)
