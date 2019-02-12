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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/spf13/cobra"
)

func init() {
	serviceManifestFlags(serviceUpdateCmd)
	serviceCmd.AddCommand(serviceUpdateCmd)
}

const serviceUpdateExample = `
  # Update info for 'redis' service in 'ns-demo' namespace
  lb service update ns-demo redis --desc "Example new description" -m 128
`

var serviceUpdateCmd = &cobra.Command{
	Use:     "update [NAMESPACE]/[NAME]",
	Short:   "Change configuration of the service",
	Example: serviceUpdateExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace, name, err := serviceParseSelfLink(args[0])
		checkError(err)

		opts, err := serviceParseManifest(cmd, name, types.EmptyString)
		checkError(err)

		cli := envs.Get().GetClient()
		response, err := cli.Cluster.V1().Namespace(namespace).Service(name).Update(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Service `%s` is updated", name))
		ss := view.FromApiServiceView(response)
		ss.Print()
	},
}
