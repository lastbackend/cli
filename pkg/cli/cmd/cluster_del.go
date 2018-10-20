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
	"errors"
	"github.com/lastbackend/cli/pkg/cli/storage"
	"github.com/spf13/cobra"
)

func init() {
	ClusterDelCmd.Flags().Bool("local", false, "Use local cluster")
	clusterCmd.AddCommand(ClusterDelCmd)
}

const clusterDelExample = `
  # Get information about cluster 
  lb cluster del name --local
`

var ClusterDelCmd = &cobra.Command{
	Use:     "del [NAME]",
	Short:   "Remove cluster",
	Example: clusterDelExample,
	Args:    cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		local, _ := cmd.Flags().GetBool("local")
		if !local {
			return errors.New("method allowed with local flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		err := storage.DelLocalCluster(name)
		if err != nil {
			panic(err)
		}

	},
}
