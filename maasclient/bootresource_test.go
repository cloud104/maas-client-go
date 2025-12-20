/*
Copyright 2021 Spectro Cloud

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package maasclient

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetBootResources(t *testing.T) {
	httpClient := &http.Client{}

	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	c := NewAuthenticatedClientSet("http://maas.test", "dummy-api-key", func(client *authenticatedClientSet) { client.WithHTTPClient(httpClient) })

	ctx := context.Background()

	t.Run("list-all", func(t *testing.T) {
		defer httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodGet,
			"http://maas.test/api/2.0/boot-resources/",
			httpmock.NewBytesResponder(200, mockData(t, "bootresources/list__all.json")),
		)

		list, err := c.BootResources().List(ctx, nil)
		assert.Nil(t, err, "expecting nil error")
		assert.NotEmpty(t, list)
	})

	t.Run("list-by-id", func(t *testing.T) {
		defer httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodGet,
			"http://maas.test/api/2.0/boot-resources/7/",
			httpmock.NewBytesResponder(200, mockData(t, "bootresources/get__id-7.json")),
		)

		res, err := c.BootResources().BootResource(7).Get(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("import image", func(t *testing.T) {
		tmp, err := os.CreateTemp("", "maas-bootresource-*.tgz")
		assert.NoError(t, err)

		defer func() { _ = os.Remove(tmp.Name()) }()

		content := []byte("dummy boot resource payload\n")
		_, err = tmp.Write(content)
		assert.NoError(t, err)
		assert.NoError(t, tmp.Close())

		size := len(content)
		sum := sha256.Sum256(content)
		sha := hex.EncodeToString(sum[:])

		defer httpmock.Reset()

		httpmock.RegisterResponder(
			http.MethodPost,
			"http://maas.test/api/2.0/boot-resources/",
			httpmock.NewStringResponder(200, fmt.Sprintf(
				string(mockData(t, "bootresources/import_response__id-99.tmpl.json")),
				size, sha, size,
			)),
		)
		httpmock.RegisterResponder(
			http.MethodPut,
			"http://maas.test/api/2.0/boot-resources/99/upload/",
			httpmock.NewStringResponder(200, `{}`),
		)

		res, err := c.BootResources().Builder("test-image",
			"amd64/generic",
			sha,
			tmp.Name(), size).Create(ctx)
		assert.Nil(t, err)
		err = res.Upload(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}
