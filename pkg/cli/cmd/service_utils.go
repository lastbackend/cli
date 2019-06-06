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

package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/spf13/cobra"
	"strings"
)

func serviceParseSelfLink(selflink string) (string, string, error) {
	match := strings.Split(selflink, "/")

	var (
		namespace, name string
	)

	switch len(match) {
	case 2:
		namespace = match[0]
		name = match[1]
	case 1:
		fmt.Println("Use default namespace:", types.DEFAULT_NAMESPACE)
		namespace = types.DEFAULT_NAMESPACE
		name = match[0]
	default:
		return "", "", errors.New("invalid service name provided")
	}

	return namespace, name, nil
}

func serviceManifestFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "set service name")
	cmd.Flags().StringP("desc", "d", "", "set service description")
	cmd.Flags().StringP("memory", "m", "128MIB", "set service spec memory")
	cmd.Flags().IntP("replicas", "r", 0, "set service replicas")
	cmd.Flags().StringArrayP("port", "p", make([]string, 0), "set service ports")
	cmd.Flags().StringArrayP("env", "e", make([]string, 0), "set service env")
	cmd.Flags().StringArray("env-from-secret", make([]string, 0), "set service env from secret")
	cmd.Flags().StringArray("env-from-config", make([]string, 0), "set service env from config")
	cmd.Flags().StringP("image", "i", "", "set service image")
	cmd.Flags().String("image-secret-name", "", "set service image auth secret name")
	cmd.Flags().String("image-secret-key", "", "set service image auth secret key")
}

func serviceParseManifest(cmd *cobra.Command, name, image string) (*request.ServiceManifest, error) {

	var err error

	description, err := cmd.Flags().GetString("desc")
	checkFlagParseError(err)

	memory, err := cmd.Flags().GetString("memory")
	checkFlagParseError(err)

	if name == types.EmptyString {
		name, err = cmd.Flags().GetString("name")
		checkFlagParseError(err)
	}

	if image == types.EmptyString {
		image, err = cmd.Flags().GetString("image")
		checkFlagParseError(err)
	}

	ports, err := cmd.Flags().GetStringArray("ports")
	checkFlagParseError(err)

	env, err := cmd.Flags().GetStringArray("env")
	checkFlagParseError(err)

	senv, err := cmd.Flags().GetStringArray("env-from-secret")
	checkFlagParseError(err)

	cenv, err := cmd.Flags().GetStringArray("env-from-config")
	checkFlagParseError(err)

	replicas, err := cmd.Flags().GetInt("replicas")
	checkFlagParseError(err)

	authName, err := cmd.Flags().GetString("image-secret-name")
	checkFlagParseError(err)

	authKey, err := cmd.Flags().GetString("image-secret-key")
	checkFlagParseError(err)

	opts := new(request.ServiceManifest)
	css := make([]request.ManifestSpecTemplateContainer, 0)

	cs := request.ManifestSpecTemplateContainer{}

	if len(name) != 0 {
		opts.Meta.Name = &name
	}

	if len(description) != 0 {
		opts.Meta.Description = &description
	}

	if memory != types.EmptyString {
		cs.Resources.Request.RAM = memory
	}

	if replicas != 0 {
		opts.Spec.Replicas = &replicas
	}

	if len(ports) > 0 {
		opts.Spec.Network = new(request.ManifestSpecNetwork)
		opts.Spec.Network.Ports = make([]string, 0)
		opts.Spec.Network.Ports = ports
	}

	es := make(map[string]request.ManifestSpecTemplateContainerEnv)
	if len(env) > 0 {
		for _, e := range env {
			kv := strings.SplitN(e, "=", 2)
			eo := request.ManifestSpecTemplateContainerEnv{
				Name: kv[0],
			}
			if len(kv) > 1 {
				eo.Value = kv[1]
			}

			es[eo.Name] = eo
		}

	}
	if len(senv) > 0 {
		for _, e := range senv {
			kv := strings.SplitN(e, "=", 3)
			eo := request.ManifestSpecTemplateContainerEnv{
				Name: kv[0],
			}
			if len(kv) < 3 {
				return nil, errors.New("Service env from secret is in wrong format, should be [NAME]=[SECRET NAME]=[SECRET STORAGE KEY]")
			}

			if len(kv) == 3 {
				eo.Secret.Name = kv[1]
				eo.Secret.Key = kv[2]
			}

			es[eo.Name] = eo
		}
	}
	if len(cenv) > 0 {
		for _, e := range cenv {
			kv := strings.SplitN(e, "=", 3)
			eo := request.ManifestSpecTemplateContainerEnv{
				Name: kv[0],
			}
			if len(kv) < 3 {
				return nil, errors.New("Service env from config is in wrong format, should be [NAME]=[CONFIG NAME]=[CONFIG KEY]")
			}

			if len(kv) == 3 {
				eo.Config.Name = kv[1]
				eo.Config.Key = kv[2]
			}

			es[eo.Name] = eo
		}
	}

	if len(es) > 0 {
		senvs := make([]request.ManifestSpecTemplateContainerEnv, 0)
		for _, e := range es {
			senvs = append(senvs, e)
		}
		cs.Env = senvs
	}

	opts.Meta.Description = &description
	cs.Image.Name = image

	if authName != types.EmptyString {
		cs.Image.Secret.Name = authName
	}

	if authKey != types.EmptyString {
		cs.Image.Secret.Key = authKey
	}

	css = append(css, cs)

	if err := opts.Validate(); err != nil {
		return nil, err.Err()
	}

	return opts, nil
}
