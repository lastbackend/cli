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
	"github.com/lastbackend/cli/pkg/client/genesis/http/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
)

type ClusterClient struct {
	client *request.RESTClient
}

func (cc *ClusterClient) List(ctx context.Context) error {

	var s *views.ClusterList
	var e *errors.Http

	err := cc.client.Get("/cluster").
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func newClusterClient(req *request.RESTClient) *ClusterClient {
	return &ClusterClient{client: req}
}
