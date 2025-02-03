package test

import (
	"os"

	"github.com/alicebob/miniredis/v2"
)

func SetupRedis() error {
	s, err := miniredis.Run()
	if err != nil {
		return err
	}
	os.Setenv("REDIS_ADDRESS", s.Addr())
	return nil
}
