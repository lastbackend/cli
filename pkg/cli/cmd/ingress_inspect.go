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
// patents in process, and are protected by trade secretCmd or copyright law.
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
	ingressCmd.AddCommand(ingressInspectCmd)
}

const ingressInspectExample = `
  # Get information 'wef34fg' for ingress
  lb ingress inspect wef34fg"
`

var ingressInspectCmd = &cobra.Command{
	Use:     "inspect [NAME]",
	Short:   "Ingress info by name",
	Example: ingressInspectExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		cli := envs.Get().GetClient()
		response, err := cli.V1().Cluster().Ingress(name).Get(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		ss := view.FromApiIngressView(response)
		ss.Print()
	},
}
