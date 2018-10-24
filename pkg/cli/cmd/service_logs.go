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
	"github.com/lastbackend/registry/pkg/distribution/types"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(serviceLogsCmd)
}

const serviceLogsExample = `
  # Get 'redis' service logs for 'ns-demo' namespace
  lb service logs ns-demo redis
`

type LogsWriter struct {
	io.Writer
}

func (LogsWriter) Write(p []byte) (int, error) {
	return fmt.Print(string(p))
}

type mapInfo map[string]serviceInfo
type serviceInfo struct {
	Deployment string
	Pod        string
	Container  string
}

var serviceLogsCmd = &cobra.Command{
	Use:     "logs [NAMESPACE] [NAME]",
	Short:   "Get service logs",
	Example: serviceLogsExample,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			choice = "0"
			m      = make(mapInfo)
			index  = 0
		)

		namespace := args[0]
		name := args[1]

		cli := envs.Get().GetClient()
		response, err := cli.Cluster.V1().Namespace(namespace).Service(name).Get(envs.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, deployment := range response.Deployments {
			state := deployment.Status.State

			if !(state == types.StateReady) {
				continue
			}

			for _, pod := range deployment.Pods {
				for _, container := range pod.Status.Containers {
					fmt.Printf("[%d] %s\n", index, container.Image.Name)
					m[strconv.Itoa(index)] = serviceInfo{
						Deployment: deployment.Meta.Name,
						Pod:        pod.Meta.Name,
						Container:  container.ID,
					}
				}
				index++
			}
		}

		if len(m) == 0 {
			fmt.Println("service in status: ", response.Status.State)
			return
		}

		for {
			fmt.Print("\nEnter container number for watch log or ^C for Exit: ")
			fmt.Scan(&choice)
			choice = strings.ToLower(choice)

			if _, ok := m[choice]; ok {
				break
			}

			fmt.Println("Number not correct!")
		}

		opts := new(request.ServiceLogsOptions)
		opts.Deployment = m[choice].Deployment
		opts.Pod = m[choice].Pod
		opts.Container = m[choice].Container

		reader, err := cli.Cluster.V1().Namespace(namespace).Service(name).Logs(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Service logs:")
		io.Copy(os.Stdout, reader)
	},
}
