package maasclient

import (
	"bytes"
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/*.json testdata/*.tmpl.json
var mockdataFS embed.FS

func mockData(t *testing.T, name string) []byte {
	t.Helper()

	b, err := mockdataFS.ReadFile("testdata/" + name)
	require.NoError(t, err)

	return bytes.TrimSpace(b)
}
