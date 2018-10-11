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
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	secretCreateCmd.Flags().StringArrayP("text", "t", make([]string, 0), "write text data in key=value format")
	secretCreateCmd.Flags().StringArrayP("file", "f", make([]string, 0), "create secret from files")
	secretCreateCmd.Flags().BoolP("auth", "a", false, "create auth secret")
	secretCreateCmd.Flags().StringP("username", "u", types.EmptyString, "add username to registry secret")
	secretCreateCmd.Flags().StringP("password", "p", types.EmptyString, "add password to registry secret")
	secretCmd.AddCommand(secretCreateCmd)
}

const secretCreateExample = `
  # Create secret 'token' with 'secret' data 
  lb secret create token secret"
`

var secretCreateCmd = &cobra.Command{
	Use:     "create [NAMESPACE] [NAME]",
	Short:   "Create secret",
	Example: secretCreateExample,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		auth, err := cmd.Flags().GetBool("auth")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		text, err := cmd.Flags().GetStringArray("text")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		files, err := cmd.Flags().GetStringArray("file")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		namespace := args[0]
		name := args[1]
		opts := new(request.SecretManifest)
		opts.Meta.Name = &name
		opts.Spec.Data = make(map[string]string, 0)

		switch true {
		case len(text) > 0:
			opts.Spec.Type = types.KindSecretOpaque

			for _, t := range text {
				var (
					k, v string
				)

				kv := strings.SplitN(t, "=", 2)
				k = kv[0]
				if len(kv) > 1 {
					v = kv[1]
				}
				opts.Spec.Data[k] = v
			}

			break
		case auth:
			opts.Spec.Type = types.KindSecretAuth
			username, err := cmd.Flags().GetString("username")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			opts.SetAuthData(username, password)
			break
		case len(files) > 0:
			opts.Spec.Type = types.KindSecretOpaque
			for _, f := range files {
				c, err := ioutil.ReadFile(f)
				if err != nil {
					_ = fmt.Errorf("failed read data from file: %s", f)
					os.Exit(1)
				}
				opts.Spec.Data[f] = string(c)
			}
			break
		default:
			fmt.Println("You need to provide secret type")
			return
		}

		if err := opts.Validate(); err != nil {
			fmt.Println(err.Err())
			return
		}

		cli := envs.Get().GetClient()
		response, err := cli.V1().Namespace(namespace).Secret().Create(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(fmt.Sprintf("Secret `%s` is created", name))

		secret := view.FromApiSecretView(response)
		secret.Print()
	},
}
