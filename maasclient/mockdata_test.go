package maasclient

import (
	"bytes"
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/*.json testdata/*.tmpl.json
var mockdataFS embed.FS

func mockData(t *testing.T, name string) string {
	t.Helper()

	b, err := mockdataFS.ReadFile("mockdata/" + name)
	require.NoError(t, err)

	return string(bytes.TrimSpace(b))
}
