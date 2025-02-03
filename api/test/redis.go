package test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func SetupRedis(t *testing.T) *miniredis.Miniredis {
	t.Helper()
	redis := miniredis.RunT(t)
	t.Setenv("REDIS_ADDRESS", redis.Addr())
	return redis
}
