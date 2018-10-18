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
	"github.com/lastbackend/cli/pkg/client/genesis/http/v1/request"
	"os"

	"github.com/howeyc/gopass"
	"github.com/lastbackend/cli/pkg/cli/config"
	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/cli/pkg/cli/storage"
	"github.com/lastbackend/cli/pkg/client"
	"github.com/lastbackend/cli/pkg/client/genesis"
	"github.com/spf13/cobra"

)

func init() {
	RootCmd.AddCommand(
		loginCmd,
		logoutCmd,
		clusterCmd,
		namespaceCmd,
		routeCmd,
		serviceCmd,
		secretCmd,
		configCmd,
		volumeCmd,
		tokenCmd,
		versionCmd,
		nodeCmd,
		ingressCmd,
	)
}

var (
	cfg = config.Get()
	ctx = envs.Get()
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lb",
	Short: "Apps cloud hosting with integrated deployment tools",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			fmt.Println("debug mode enabled")
			cfg.Debug = debug
		}

		token, err := storage.GetToken()
		if err != nil {
			panic("There is no token in .lastbackend in homedir")
		}

		//host := cmd.Flag("host").Value.String()
		//if len(host) == 0 {
		host := "https://api.lastbackend.com"
		//}

		gcfg := genesis.NewConfig()
		gcfg.BearerToken = token

		gccli, err := genesis.New(genesis.ClientHTTP, host, gcfg)
		if err != nil {
			panic(err)
		}

		//rcfg := registry.NewConfig()
		//rcfg.BearerToken = token
		//
		//cccli, err := cluster.New(cluster.ClientHTTP, host, rcfg)
		//if err != nil {
		//	panic(err)
		//}
		//
		//ccfg := cluster.NewConfig()
		//ccfg.BearerToken = token

		//rccli, err := registry.New(registry.ClientHTTP, host, ccfg)
		//if err != nil {
		//	panic(err)
		//}

		cli := &client.Client{
			//Cluster:  cccli,
			//Registry: rccli,
			Genesis: gccli,
		}

		ctx.SetClient(cli)
	},
}

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage your namespaces",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		command := "[string]"
		if len(args) != 0 {
			command = args[0]
		}

		var ns = &cobra.Command{
			Use:   command,
			Short: "Manage your a namespace",
		}

		cmd.AddCommand(ns)

		if len(args) == 0 {
			cmd.Help()
			return
		}

		// Attach sub command for namespace
		ns.AddCommand(
			serviceCmd,
			secretCmd,
			routeCmd,
		)

		ns.Execute()

	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a Last.Backend",
	Example: `
  # Log in to a Last.Backend 
  lb login
  Login: username
  Password: ******"`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			login    string
			password string
		)

		fmt.Print("Login: ")
		fmt.Scan(&login)

		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println(err)
			return
		}
		password = string(pass)

		fmt.Print("\r\n")

		cli := envs.Get().GetClient()

		opts := &request.AccountLoginOptions{
			Login:    login,
			Password: password,
		}

		session, err := cli.Genesis.V1().Account().Login(envs.Background(), opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := storage.SetToken(session.Token); err != nil {
			fmt.Println(err)
			return
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from a Last.Backend",
	Example: `
  # Log out from a Last.Backend 
  lb logout"`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := storage.SetToken(""); err != nil {
			fmt.Println(err)
			return
		}
	},
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage your service",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage your secret",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage your configs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Manage your volumes",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage your route",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage your cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage set vars to your local storage",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage cluster nodes",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "Manage cluster ingress servers",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Manage cluster discovery servers",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command
func Execute() {

	cobra.OnInitialize()

	RootCmd.PersistentFlags().StringP("host", "H", "https://api.lastbackend.com", "Set api host parameter")
	RootCmd.PersistentFlags().Bool("local", false, "Set api host for local cluster")
	RootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	RootCmd.PersistentFlags().Bool("insecure", false, "Disable security check")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
