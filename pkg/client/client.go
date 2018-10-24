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
	gc "github.com/lastbackend/cli/pkg/client/genesis"
	cc "github.com/lastbackend/lastbackend/pkg/api/client"
	rc "github.com/lastbackend/registry/pkg/api/client"
)

type Client struct {
	Genesis  gc.IClient
	Cluster  cc.IClient
	Registry rc.IClient
}

func NewGenesisClister(host string, config *Config) gc.IClient {

	if config == nil {
		config = new(Config)
	}

	cfg := gc.NewConfig()

	cfg.BearerToken = config.Token

	if config.TLS != nil {
		cfg.TLS.CertFile = config.TLS.CAFile
		cfg.TLS.KeyFile = config.TLS.KeyFile
		cfg.TLS.CAFile = config.TLS.CAFile
		cfg.TLS.CertData = config.TLS.CertData
		cfg.TLS.KeyData = config.TLS.KeyData
		cfg.TLS.CAData = config.TLS.CAData
	}

	cli, err := gc.New(gc.ClientHTTP, host, cfg)
	if err != nil {
		panic(err)
	}

	return cli
}

func NewClusterClient(host string, config *Config) cc.IClient {

	if config == nil {
		config = new(Config)
	}

	cfg := cc.NewConfig()

	if config.Headers != nil {
		cfg.Headers = make(map[string]string, 0)
		for k, v := range config.Headers {
			cfg.Headers[k] = v
		}
	}

	cfg.BearerToken = config.Token

	if config.TLS != nil {
		cfg.TLS.CertFile = config.TLS.CAFile
		cfg.TLS.KeyFile = config.TLS.KeyFile
		cfg.TLS.CAFile = config.TLS.CAFile
		cfg.TLS.CertData = config.TLS.CertData
		cfg.TLS.KeyData = config.TLS.KeyData
		cfg.TLS.CAData = config.TLS.CAData
	}

	cli, err := cc.New(cc.ClientHTTP, host, cfg)
	if err != nil {
		panic(err)
	}

	return cli
}

func NewRegistryClient(host string, config *Config) rc.IClient {

	if config == nil {
		config = new(Config)
	}

	cfg := rc.NewConfig()

	cfg.BearerToken = config.Token

	if config.TLS != nil {
		cfg.TLS.CertFile = config.TLS.CAFile
		cfg.TLS.KeyFile = config.TLS.KeyFile
		cfg.TLS.CAFile = config.TLS.CAFile
		cfg.TLS.CertData = config.TLS.CertData
		cfg.TLS.KeyData = config.TLS.KeyData
		cfg.TLS.CAData = config.TLS.CAData
	}

	cli, err := rc.New(cc.ClientHTTP, host, cfg)
	if err != nil {
		panic(err)
	}

	return cli
}
