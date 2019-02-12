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
	"github.com/spf13/cobra"
)

func init() {
	jobCmd.AddCommand(jobRunCmd)
}

const jobRunExample = `
  # Get information for 'redis' job in 'ns-demo' namespace
  lb job run ns/cron
`

var jobRunCmd = &cobra.Command{
	Use:     "run [NAMESPACE]/[NAME]",
	Short:   "Run job info by name",
	Example: jobRunExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace, name, err := jobParseSelfLink(args[0])
		checkError(err)

		cli := envs.Get().GetClient()
		_, err = cli.Cluster.V1().Namespace(namespace).Job(name).Run(envs.Background(), nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
