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

			item, err := storage.GetLocalCluster(name)
			if err != nil {
				fmt.Println(err)
			}

			if item == nil {
				fmt.Println(fmt.Sprintf("Cluster `%s` not found", name))
			}

			err = storage.SetCluster(fmt.Sprintf("l.%s", item.Name))
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(fmt.Sprintf("Cluster `%s` selected", name))

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

		err = storage.SetCluster(fmt.Sprintf("r.%s", cl.Meta.SelfLink))
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("Cluster `%s` selected", name))
	},
}
