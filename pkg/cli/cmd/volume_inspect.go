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
	"github.com/spf13/cobra"
)

func init() {
	volumeCmd.AddCommand(volumeInspectCmd)
}

const volumeInspectExample = `
  # Get information for 'redis' volume in 'ns-demo' namespace
  lb volume inspect ns-demo redis
`

var volumeInspectCmd = &cobra.Command{
	Use:     "inspect [NAMESPACE] [NAME]",
	Short:   "Volume info by name",
	Example: volumeInspectExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		cli := envs.Get().GetClient()
		response, err := cli.Cluster.V1().Namespace(namespace).Volume(name).Get(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		ss := view.FromApiVolumeView(response)
		ss.Print()
	},
}
