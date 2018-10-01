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
	discoveryCmd.AddCommand(discoveryListCmd)
}

const discoveryListExample = `
  # Get all discoveryes in cluster 
  lb discovery ls
`

var discoveryListCmd = &cobra.Command{
	Use:     "ls",
	Short:   "Display the discoverys list",
	Example: discoveryListExample,
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		cli := envs.Get().GetClient()
		response, err := cli.V1().Cluster().Discovery().List(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		if response == nil || len(*response) == 0 {
			fmt.Println("no discoverys available")
			return
		}

		list := view.FromApiDiscoveryListView(response)
		list.Print()
	},
}
