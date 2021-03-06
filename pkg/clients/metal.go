/*
Copyright 2019 The Crossplane Authors.

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

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/packethost/packngo"

	"github.com/packethost/crossplane-provider-equinix-metal/pkg/version"
)

// Client is a structure that embeds Credentials for the purposes of
// defaulting to those credential supplied values during Equinix Metal API usage. This
// allows for the Device resource to not require a ProjectID, for example, since
// the provider was configured with a ProjectID.
type Client struct {
	*Credentials

	Client *packngo.Client
}

// NewClient returns an Equinix Metal Client configured with credentials
func NewClient(ctx context.Context, credentials []byte) (*Client, error) {
	config := &Credentials{}
	if err := json.Unmarshal(credentials, config); err != nil {
		return nil, err
	}

	apiKey := config.GetAPIKey(CredentialAPIKey)
	if apiKey == "" {
		return nil, fmt.Errorf("Invalid APIKey in credentials")
	}
	apiClient := packngo.NewClientWithAuth("crossplane", apiKey, nil)
	apiClient.UserAgent = fmt.Sprintf("crossplane-provider-equinix-metal/%s %s", version.Version, apiClient.UserAgent)

	client := &Client{
		Client:      apiClient,
		Credentials: config,
	}

	return client, nil
}

// IsNotFound returns true if error is not found
func IsNotFound(err error) bool {
	if e, ok := err.(*packngo.ErrorResponse); ok {
		return e.Response.StatusCode == http.StatusNotFound
	}
	return false
}
