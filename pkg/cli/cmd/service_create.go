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
// patents in process, and are protected by trade secretCmd or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package cmd

import (
	"fmt"

	"strings"

	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/cli/pkg/cli/view"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/spf13/cobra"
)

func init() {
	serviceCreateCmd.Flags().StringP("file", "f", "", "create service from file")
	serviceCreateCmd.Flags().StringP("desc", "d", "", "set service description")
	serviceCreateCmd.Flags().StringP("name", "n", "", "set service name")
	serviceCreateCmd.Flags().String("image-secret", "", "service image auth secret")
	serviceCreateCmd.Flags().Int64P("memory", "m", 128, "set service spec memory")
	serviceCreateCmd.Flags().IntP("replicas", "r", 1, "set service replicas")
	serviceCreateCmd.Flags().StringArrayP("port", "p", make([]string, 0), "set service ports")
	serviceCreateCmd.Flags().StringArrayP("env", "e", make([]string, 0), "set service env")
	serviceCreateCmd.Flags().StringArray("env-from-secret", make([]string, 0), "set service env from secret")
	serviceCmd.AddCommand(serviceCreateCmd)
}

const serviceCreateExample = `
  # Create new redis service with description and 256 MB limit memory
  lb service create ns-demo redis --desc "Example description" -m 256
`

var serviceCreateCmd = &cobra.Command{
	Use:     "create [NAMESPACE] [IMAGE]",
	Short:   "Create service",
	Example: serviceCreateExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		image := args[1]

		description, _ := cmd.Flags().GetString("desc")
		memory, _ := cmd.Flags().GetInt64("memory")
		name, _ := cmd.Flags().GetString("name")
		ports, _ := cmd.Flags().GetStringArray("ports")
		env, _ := cmd.Flags().GetStringArray("env")
		senv, _ := cmd.Flags().GetStringArray("env-from-secret")
		replicas, _ := cmd.Flags().GetInt("replicas")
		auth, _ := cmd.Flags().GetString("image-secret")

		opts := new(request.ServiceManifest)
		css := make([]request.ManifestSpecTemplateContainer, 0)

		cs := request.ManifestSpecTemplateContainer{}

		if len(name) != 0 {
			opts.Meta.Name = &name
		}

		if len(description) != 0 {
			opts.Meta.Description = &description
		}

		if memory != 0 {
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
					fmt.Println("Service env from secret is in wrong format, should be [NAME]=[SECRET NAME]=[SECRET STORAGE KEY]")
					return
				}

				if len(kv) == 3 {
					eo.From.Name = kv[1]
					eo.From.Key = kv[2]
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

		if auth != types.EmptyString {
			cs.Image.Secret = auth
		}

		css = append(css, cs)

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Service().Create(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Service `%s` is created", name))

		service := view.FromApiServiceView(response)
		service.Print()
	},
}
