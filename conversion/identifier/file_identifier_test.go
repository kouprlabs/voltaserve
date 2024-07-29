package identifier_test

import (
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestFileIdentifier_IsKRA(t *testing.T) {
	tests := map[string]struct {
		FileName     string
		IsKRA        bool
		RequireError require.ErrorAssertionFunc
	}{
		"KritaFile":    {FileName: "krita.kra", IsKRA: true, RequireError: require.NoError},
		"NormalZip":    {FileName: "regular.zip", IsKRA: false, RequireError: require.NoError},
		"MimeTypedZip": {FileName: "mimetyped.zip", IsKRA: false, RequireError: require.NoError},
		"NoSuchFile":   {FileName: "doesnotexist", IsKRA: false, RequireError: require.Error},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fi := identifier.NewFileIdentifier()

			isKra, err := fi.IsKRA(filepath.Join("fixtures", tc.FileName))
			tc.RequireError(t, err, "IsKRA(%#v)", tc.FileName)

			assert.Equal(t, tc.IsKRA, isKra, "IsKRA(%#v)", tc.FileName)
		})
	}
}
