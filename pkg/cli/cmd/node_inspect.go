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
	nodeCmd.AddCommand(nodeInspectCmd)
}

const nodeInspectExample = `
  # Get information 'wef34fg' for node in 'ns-demo' namespace
  lb node inspect ns-demo wef34fg"
`

var nodeInspectCmd = &cobra.Command{
	Use:     "inspect [NAME]",
	Short:   "Node info by name",
	Example: nodeInspectExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		cli := envs.Get().GetClient()
		response, err := cli.Cluster.V1().Cluster().Node(name).Get(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		ss := view.FromApiNodeView(response)
		ss.Print()
	},
}
