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
	secretCmd.AddCommand(secretInspectCmd)
}

const secretInspectExample = `
  # Inspect secret 'token' 
  lb secret inspect token"
`

var secretInspectCmd = &cobra.Command{
	Use:     "inspect [NAMESPACE] [NAME]",
	Short:   "Inspect secret",
	Example: secretInspectExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]
		name := args[1]

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Secret(name).Get(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		secret := view.FromApiSecretView(response)
		secret.Print()
	},
}
