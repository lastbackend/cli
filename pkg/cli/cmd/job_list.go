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
	jobCmd.AddCommand(jobListCmd)
}

const jobListExample = `
  # Get all jobs for 'ns-demo' namespace  
  lb job ls ns-demo
`

var jobListCmd = &cobra.Command{
	Use:     "ls [NAMESPACE]",
	Short:   "Display the jobs list",
	Example: jobListExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]

		cli := envs.Get().GetClient()

		response, err := cli.Cluster.V1().Namespace(namespace).Job().List(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		if response == nil || len(*response) == 0 {
			fmt.Println("no jobs available")
			return
		}

		list := view.FromApiJobListView(response)
		list.Print()
	},
}
