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
	"github.com/lastbackend/cli/pkg/cli/view"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/decoder"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const applyExample = `
  # Apply manifest from file or by URL
  lb namespace [name] apply -f"
`

func init() {
	applyCmd.Flags().StringArrayP("file", "f", make([]string, 0), "apply resources to namespace from files")
	namespaceCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply [NAME]",
	Short: "Apply manifest files to cluster",
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
		spec := v1.Request().Namespace().ApplyManifest()

		for _, f := range files {

			s, err := os.Open(f)
			if err != nil {
				if os.IsNotExist(err) {
					_ = fmt.Errorf("failed read data: file not exists: %s", f)
					os.Exit(1)
				}
			}

			if err := s.Close(); err != nil {
				fmt.Errorf("close file err: %s", err.Error())
				return
			}

			c, err := ioutil.ReadFile(f)
			if err != nil {
				_ = fmt.Errorf("failed read data from file: %s", f)
				os.Exit(1)
			}

			items := decoder.YamlSplit(c)
			fmt.Println("manifests:", len(items))

			for _, i := range items {

				var m = new(request.Runtime)

				if err := yaml.Unmarshal([]byte(i), m); err != nil {
					_ = fmt.Errorf("can not parse manifest: %s: %s", f, err.Error())
					continue
				}

				switch strings.ToLower(m.Kind) {
				case types.KindConfig:
					m := new(request.ConfigManifest)
					err := m.FromYaml(i)
					if err != nil {
						_ = fmt.Errorf("invalid specification: %s", err.Error())
						return
					}
					if m.Meta.Name == nil {
						break
					}
					fmt.Printf("Add config manifest: %s\n", *m.Meta.Name)
					spec.Configs[*m.Meta.Name] = m
					break
				case types.KindSecret:
					m := new(request.SecretManifest)
					err := m.FromYaml(i)
					if err != nil {
						_ = fmt.Errorf("invalid specification: %s", err.Error())
						return
					}
					if m.Meta.Name == nil {
						break
					}
					fmt.Printf("Add secret manifest: %s\n", *m.Meta.Name)
					spec.Secrets[*m.Meta.Name] = m
					break
				case types.KindService:
					m := new(request.ServiceManifest)
					err := m.FromYaml(i)
					if err != nil {
						_ = fmt.Errorf("invalid specification: %s", err.Error())
						return
					}
					if m.Meta.Name == nil {
						break
					}
					fmt.Printf("Add service manifest: %s\n", *m.Meta.Name)
					spec.Services[*m.Meta.Name] = m
					break
				case types.KindVolume:

					m := new(request.VolumeManifest)
					err := m.FromYaml(i)
					if err != nil {
						_ = fmt.Errorf("invalid specification: %s", err.Error())
						return
					}
					if m.Meta.Name == nil {
						break
					}
					fmt.Printf("Add volume manifest: %s\n", *m.Meta.Name)
					spec.Volumes[*m.Meta.Name] = m
					break
				case types.KindJob:
					m := new(request.JobManifest)
					err := m.FromYaml(i)
					if err != nil {
						_ = fmt.Errorf("invalid specification: %s", err.Error())
						return
					}
					if m.Meta.Name == nil {
						break
					}
					fmt.Printf("Add job manifest: %s\n", *m.Meta.Name)
					spec.Jobs[*m.Meta.Name] = m
					break
				case types.KindRoute:
					m := new(request.RouteManifest)
					err := m.FromYaml(i)
					if err != nil {
						_ = fmt.Errorf("invalid specification: %s", err.Error())
						return
					}
					if m.Meta.Name == nil {
						break
					}
					fmt.Printf("Add route manifest: %s\n", *m.Meta.Name)
					spec.Routes[*m.Meta.Name] = m
					break
				}
			}

			status, err := cli.Cluster.V1().Namespace(namespace).Apply(envs.Background(), spec)
			if err != nil {
				_ = fmt.Errorf("invalid specification: %s", err.Error())
				return
			}

			fmt.Println()
			ns := view.FromApiNamespaceStatusView(status)
			ns.Print()
			return

		}
	},
}
