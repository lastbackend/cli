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
	"github.com/lastbackend/cli/pkg/cli/storage"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	ClusterAddCmd.Flags().Bool("local", false, "Use local cluster")
	clusterCmd.AddCommand(ClusterAddCmd)
}

const clusterAddExample = `
  # Get information about cluster 
  lb cluster add name endpoint --local
`

var ClusterAddCmd = &cobra.Command{
	Use:     "add [NAME] [ENDPOINT]",
	Short:   "Add cluster",
	Example: clusterAddExample,
	Args:    cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		local, _ := cmd.Flags().GetBool("local")
		if !local {
			return errors.New("method allowed with local flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		endpoint := args[1]

		local, err := cmd.Flags().GetBool("local")
		if err != nil {
			panic(err)
		}

		err = storage.AddLocalCluster(name, endpoint, local)
		switch true {
		case err == nil:
		case err.Error() == "already exists":
			fmt.Println(fmt.Sprintf("Cluster `%s` already exists", name))
		default:
			panic(err)
		}

	},
}
