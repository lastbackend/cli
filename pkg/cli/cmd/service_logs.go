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
	"encoding/json"
	"fmt"
	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	serviceCmd.AddCommand(serviceLogsCmd)
}

const serviceLogsExample = `
  # Get 'redis' service logs for 'ns-demo' namespace
  lb service logs ns-demo redis
`

var serviceLogsCmd = &cobra.Command{
	Use:     "logs [NAMESPACE]/[NAME]",
	Short:   "Get service logs",
	Example: serviceLogsExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		opts := new(request.ServiceLogsOptions)

		namespace, name, err := serviceParseSelfLink(args[0])
		checkError(err)

		cli := envs.Get().GetClient()

		reader, err := cli.Cluster.V1().Namespace(namespace).Service(name).Logs(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		dec := json.NewDecoder(reader)
		for {
			var doc types.LogMessage

			err := dec.Decode(&doc)
			if err == io.EOF {
				// all done
				break
			}
			if err != nil {
				fmt.Errorf(err.Error())
				os.Exit(1)
			}

			fmt.Println(">", doc.Selflink, doc.Data)
		}
	},
}
