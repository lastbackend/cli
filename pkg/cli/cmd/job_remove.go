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
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobRemoveCmd)
}

const jobRemoveExample = `
  # Remove 'redis' job in 'ns-demo' namespace
  lb job remove ns-demo redis
`

var jobRemoveCmd = &cobra.Command{
	Use:     "remove [NAMESPACE] [NAME]",
	Short:   "Remove job by name",
	Example: jobRemoveExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		opts := &request.JobRemoveOptions{Force: false}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		if err := cli.Cluster.V1().Namespace(namespace).Job(name).Remove(envs.Background(), opts); err != nil {
			_ = fmt.Errorf("job remove err: %s", err.Error())
			return
		}

		fmt.Println(fmt.Sprintf("Job `%s` is successfully removed", name))
	},
}
