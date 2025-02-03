package cache_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kouprlabs/voltaserve/api/test"
)

func TestMain(m *testing.M) {
	postgres, err := test.SetupPostgres(25432)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := postgres.Stop(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()
	if err := test.SetupRedis(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}
