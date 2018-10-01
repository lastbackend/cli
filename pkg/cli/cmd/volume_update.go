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
	"github.com/spf13/cobra"
)

func init() {
	volumeUpdateCmd.Flags().StringP("desc", "d", "", "set volume description")
	volumeCmd.AddCommand(volumeUpdateCmd)
}

const volumeUpdateExample = `
  # Update info for 'redis' volume in 'ns-demo' namespace
  lb volume update ns-demo redis --desc "Example new description" -m 128
`

var volumeUpdateCmd = &cobra.Command{
	Use:     "update [NAMESPACE] [NAME]",
	Short:   "Change configuration of the volume",
	Example: volumeUpdateExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		description, _ := cmd.Flags().GetString("desc")

		opts := new(request.VolumeManifest)

		if len(name) != 0 {
			opts.Meta.Name = &name
		}

		if len(description) != 0 {
			opts.Meta.Description = &description
		}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Volume(name).Update(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Volume `%s` is updated", name))
		ss := view.FromApiVolumeView(response)
		ss.Print()
	},
}
