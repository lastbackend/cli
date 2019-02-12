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
	jobInspectCmd.Flags().StringP("task", "t", "", "inspect particular task")
	jobCmd.AddCommand(jobInspectCmd)
}

const jobInspectExample = `
  # Get information for 'redis' job in 'ns-demo' namespace
  lb job inspect ns-demo redis
`

var jobInspectCmd = &cobra.Command{
	Use:     "inspect [NAMESPACE]/[NAME]",
	Short:   "Service info by name",
	Example: jobInspectExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace, name, err := jobParseSelfLink(args[0])
		checkError(err)

		cli := envs.Get().GetClient()

		t, err := cmd.Flags().GetString("task")
		if err != nil {
			_ = fmt.Errorf("can not parse task option: %s", t)
			return
		}

		if t == "" {
			job, err := cli.Cluster.V1().Namespace(namespace).Job(name).Get(envs.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			ss := view.FromApiJobView(job)
			ss.Print()
			return
		}

		task, err := cli.Cluster.V1().Namespace(namespace).Job(name).Tasks(t).Get(envs.Background())

		tw := view.FromApiTaskView(task)
		tw.Print()

	},
}
