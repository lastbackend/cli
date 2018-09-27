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
	"io/ioutil"
	"os"

	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/cli/pkg/cli/view"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const applyExample = `
  # Apply manifest from file or by URL
  lb namespace [name] apply -f"
`

func init() {
	applyCmd.Flags().StringArrayP("file", "f", make([]string, 0), "create secret from files")
	namespaceCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply [NAME]",
	Short: "Apply file manifest to cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		namespace := args[0]

		files, err := cmd.Flags().GetStringArray("file")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(files) == 0 {
			cmd.Help()
		}

		cli := envs.Get().GetClient()

		for _, f := range files {

			s, err := os.Open(f)
			if err != nil {
				if os.IsNotExist(err) {
					_ = fmt.Errorf("failed read data: file not exists: %s", f)
					os.Exit(1)
				}
			}
			s.Close()

			c, err := ioutil.ReadFile(f)
			if err != nil {
				_ = fmt.Errorf("failed read data from file: %s", f)
				os.Exit(1)
			}

			var m = new(request.Runtime)
			yaml.Unmarshal(c, m)

			if m.Kind == "Service" {

				spec := v1.Request().Service().Manifest()
				err := spec.FromYaml(c)
				if err != nil {
					_ = fmt.Errorf("invalid specification: %s", err.Error())
					return
				}

				var rsvc *views.Service

				if spec.Meta.Name != nil {
					rsvc, _ = cli.V1().Namespace(namespace).Service(*spec.Meta.Name).Get(envs.Background())
				}

				if rsvc == nil {
					fmt.Println("create new service")
					rsvc, err = cli.V1().Namespace(namespace).Service().Create(envs.Background(), spec)
					if err != nil {
						fmt.Println(err)
						return
					}
				} else {
					rsvc, err = cli.V1().Namespace(namespace).Service(rsvc.Meta.Name).Update(envs.Background(), spec)
					if err != nil {
						fmt.Println(3)
						fmt.Println(err)
						return
					}
				}

				if rsvc != nil {
					service := view.FromApiServiceView(rsvc)
					service.Print()
				} else {
					fmt.Println("ooops")
				}

			}

			if m.Kind == "Route" {
				spec := v1.Request().Route().Manifest()
				err := spec.FromYaml(c)
				if err != nil {
					fmt.Errorf("invalid specification: %s", err.Error())
					return
				}

				var rr *views.Route

				if spec.Meta.Name != nil {
					rr, _ = cli.V1().Namespace(namespace).Route(*spec.Meta.Name).Get(envs.Background())
				}

				if rr == nil {
					fmt.Println("create new route")
					rr, err = cli.V1().Namespace(namespace).Route().Create(envs.Background(), spec)
					if err != nil {
						fmt.Println(err)
						return
					}
				} else {
					fmt.Println("update route")
					rr, err = cli.V1().Namespace(namespace).Route(rr.Meta.Name).Update(envs.Background(), spec)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if rr != nil {
					route := view.FromApiRouteView(rr)
					route.Print()
				} else {
					fmt.Println("ooops")
				}
			}

			if m.Kind == "Volume" {
				spec := v1.Request().Volume().Manifest()
				err := spec.FromYaml(c)
				if err != nil {
					fmt.Errorf("invalid specification: %s", err.Error())
					return
				}

				var rr *views.Volume

				if spec.Meta.Name != nil {
					rr, _ = cli.V1().Namespace(namespace).Volume(*spec.Meta.Name).Get(envs.Background())
				}

				if rr == nil {
					fmt.Println("create new route")
					rr, err = cli.V1().Namespace(namespace).Volume().Create(envs.Background(), spec)
					if err != nil {
						fmt.Println(err)
						return
					}
				} else {
					fmt.Println("update route")
					rr, err = cli.V1().Namespace(namespace).Volume(rr.Meta.Name).Update(envs.Background(), spec)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if rr != nil {
					route := view.FromApiVolumeView(rr)
					route.Print()
				} else {
					fmt.Println("ooops")
				}
			}

			return

		}
	},
}
