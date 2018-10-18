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

package envs

import (
	"context"
	"github.com/lastbackend/cli/pkg/client"

	"github.com/lastbackend/cli/pkg/cli/config"
)

var _ctx ctx

func Get() *ctx {
	return &_ctx
}

func Background() context.Context {
	return context.Background()
}

func Mock() *ctx {
	_ctx.mock = true
	return &_ctx
}

type ctx struct {
	config *config.Config
	client *client.Client
	token  *string
	mock   bool
}

func (c *ctx) SetClient(client *client.Client) {
	c.client = client
}

func (c *ctx) GetClient() *client.Client {
	return c.client
}

func (c *ctx) SetSessionToken(token string) {
	c.token = &token
}

func (c *ctx) GetSessionToken() *string {
	return c.token
}

func (c *ctx) SetConfig(cfg *config.Config) {
	c.config = cfg
}

func (c *ctx) GetConfig() *config.Config {
	return c.config
}
