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
	serviceManifestFlags(serviceCreateCmd)
	serviceCmd.AddCommand(serviceCreateCmd)
}

const serviceCreateExample = `
  # Create new redis service with description and 256 MB limit memory
  lb service create ns-demo redis --desc "Example description" -m 256mib
`

var serviceCreateCmd = &cobra.Command{
	Use:     "create [NAMESPACE]/[NAME] [IMAGE]",
	Short:   "Create service",
	Example: serviceCreateExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace, name, err := serviceParseSelfLink(args[0])
		checkError(err)

		image := args[1]

		opts, err := serviceParseManifest(cmd, name, image)
		checkError(err)

		cli := envs.Get().GetClient()
		response, err := cli.Cluster.V1().Namespace(namespace).Service().Create(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Service `%s` is created", name))

		service := view.FromApiServiceView(response)
		service.Print()
	},
}
