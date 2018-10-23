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
	"github.com/lastbackend/cli/pkg/cli/storage"
	"github.com/spf13/cobra"
)

func init() {
	ClusterSelectCmd.Flags().Bool("local", false, "Use local cluster")
	clusterCmd.AddCommand(ClusterSelectCmd)
}

const clusterSelectExample = `
  # Get information about cluster 
  lb cluster select name
`

var ClusterSelectCmd = &cobra.Command{
	Use:     "select [NAME]",
	Short:   "Select cluster",
	Example: clusterSelectExample,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		// Select local cluster
		local, _ := cmd.Flags().GetBool("local")
		if local {
			items, err := storage.ListLocalCluster()
			if err != nil {
				panic(err)
			}

			if e, ok := items[name]; !ok {
				fmt.Println(fmt.Sprintf("Cluster `%s` not found", name))
			} else {
				err = storage.SetCluster(name, e, true)
				if err != nil {
					panic(err)
				}
			}

			return
		}

		// Select remove cluster
		cli := envs.Get().GetClient()

		cl, err := cli.Genesis.V1().Cluster().Get(envs.Background(), name)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cl == nil {
			fmt.Println(fmt.Sprintf("Cluster `%s` not found", name))
			return
		}

		err = storage.SetCluster(name, "", false)
		if err != nil {
			panic(err)
		}
	},
}
