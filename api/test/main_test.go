package test

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	postgres, err := setupPostgres(15432)
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
	if err := setupRedis(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}
