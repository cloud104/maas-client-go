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
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSSHKeys(t *testing.T) {
	httpClient := &http.Client{}

	httpmock.ActivateNonDefault(httpClient)
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder(
		http.MethodGet,
		"http://maas.test/api/2.0/account/prefs/sshkeys/",
		httpmock.NewStringResponder(200, `[
			{
				"id": 1,
				"key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCtestkeymaterial test@example",
				"comment": "test@example",
				"resource_uri": "/MAAS/api/2.0/account/prefs/sshkeys/1/"
			}
		]`),
	)

	c := NewAuthenticatedClientSet("http://maas.test", "dummy-api-key", func(client *authenticatedClientSet) { client.WithHTTPClient(httpClient) })

	ctx := context.Background()

	t.Run("list sshkeys", func(t *testing.T) {
		sshKeys, err := c.SSHKeys().List(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, sshKeys)
		assert.NotEmpty(t, sshKeys)
	})
}
