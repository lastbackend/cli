//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package client

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
)

type RegistryClient struct {
	client *request.RESTClient
}

func (ac *RegistryClient) Get(ctx context.Context) error {

	return nil
}

func newRegistryClient(req *request.RESTClient) *RegistryClient {
	return &RegistryClient{client: req}
}
