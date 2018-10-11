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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"

	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/cli/pkg/cli/view"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

func init() {
	volumeCreateCmd.Flags().StringP("type", "t", "", "set volume type")
	volumeCreateCmd.Flags().StringP("desc", "d", "", "set volume description")
	volumeCmd.AddCommand(volumeCreateCmd)
}

const volumeCreateExample = `
  # Create new redis volume with description and 256 MB limit memory
  lb volume create ns-demo redis --desc "Example description" -m 256
`

var volumeCreateCmd = &cobra.Command{
	Use:     "create [NAMESPACE] [NAME]",
	Short:   "Create volume",
	Example: volumeCreateExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		description, _ := cmd.Flags().GetString("desc")
		kind, _ := cmd.Flags().GetString("type")

		opts := new(request.VolumeManifest)

		if len(name) != 0 {
			opts.Meta.Name = &name
		}

		if len(description) != 0 {
			opts.Meta.Description = &description
		}

		switch kind {
		case types.KindVolumeHostDir:
			opts.Spec.Type = types.KindVolumeHostDir
			break
		}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Volume().Create(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Volume `%s` is created", name))

		volume := view.FromApiVolumeView(response)
		volume.Print()
	},
}
