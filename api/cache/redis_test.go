package cache_test

import (
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

//nolint:paralleltest
func TestRedis(t *testing.T) {
	fileCache := cache.NewFileCache()
	opts := repo.NewFileOptions{
		ID:   helper.NewID(),
		Name: "file",
		Type: model.FileTypeFile,
	}
	if err := fileCache.Set(repo.NewFileWithOptions(opts)); err != nil {
		t.Fatal(err)
	}
	file, err := fileCache.Get(opts.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, file.GetID(), opts.ID)
	assert.Equal(t, file.GetName(), opts.Name)
	assert.Equal(t, file.GetType(), opts.Type)
}
