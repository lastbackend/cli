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
	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/cli/pkg/cli/view"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	configUpdateCmd.Flags().StringArrayP("text", "t", make([]string, 0), "write config in key=value format")
	configUpdateCmd.Flags().StringArrayP("file", "f", make([]string, 0), "create config from files")
	configCmd.AddCommand(configUpdateCmd)
}

const configUpdateExample = `
  # Update 'token' config record with 'new-config' data
  lb config update token new-config"
`

var configUpdateCmd = &cobra.Command{
	Use:     "update [NAME]",
	Short:   "Change configuration of the config",
	Example: configUpdateExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		text, _ := cmd.Flags().GetStringArray("text")
		files, _ := cmd.Flags().GetStringArray("file")

		namespace := args[0]
		name := args[1]
		opts := new(request.ConfigManifest)
		opts.Spec.Data = make(map[string]string, 0)

		switch true {
		case len(text) > 0:
			opts.Spec.Type = types.KindConfigText

			for _, t := range text {
				var (
					k, v string
				)

				kv := strings.SplitN(t, "=", 2)
				k = kv[0]
				if len(kv) > 1 {
					v = kv[1]
				}
				opts.Spec.Data[k] = v
			}

			break
		case len(files) > 0:
			opts.Spec.Type = types.KindConfigText
			for _, f := range files {
				c, err := ioutil.ReadFile(f)
				if err != nil {
					_ = fmt.Errorf("failed read data from file: %s", f)
					os.Exit(1)
				}
				opts.Spec.Data[f] = string(c)
			}
			break
		default:
			fmt.Println("You need to provide config type")
			return
		}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.Cluster.V1().Namespace(namespace).Config(name).Update(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Config `%s` is updated", name))
		ss := view.FromApiConfigView(response)
		ss.Print()
	},
}
