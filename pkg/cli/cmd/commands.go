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
	"github.com/lastbackend/lastbackend/pkg/log"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/lastbackend/cli/pkg/cli/config"
	"github.com/lastbackend/cli/pkg/cli/envs"
	"github.com/lastbackend/cli/pkg/cli/storage"
	"github.com/lastbackend/cli/pkg/client"
	"github.com/lastbackend/cli/pkg/client/genesis/http/v1/request"
	"github.com/lastbackend/cli/pkg/util/filesystem"
	"github.com/spf13/cobra"
)

const defaultHost = "https://api.lastbackend.com"

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

		var err error

		cfg.Debug, err = cmd.Flags().GetBool("debug")
		if err != nil {
			fmt.Println("Invalid debug flag")
			return
		}

		token, err := storage.GetToken()
		if err != nil {
			fmt.Println("There is no token in .lastbackend in homedir")
			return
		}

		host := defaultHost
		config := &client.Config{Token: token}

		tls, err := cmd.Flags().GetBool("tls")
		if err != nil {
			fmt.Println("Invalid tls flag")
			return
		}

		if tls {
			config.TLS = new(client.TLSConfig)
			config.TLS.Insecure = false
			config.TLS.CAFile = cmd.Flag("tlscacert").Value.String()
			config.TLS.CertFile = cmd.Flag("tlscert").Value.String()
			config.TLS.KeyFile = cmd.Flag("tlskey").Value.String()
		}

		cli := &client.Client{}
		cli.Genesis = client.NewGenesisClister(host, config)
		cli.Registry = client.NewRegistryClient(host, config)

		// ============================
		// Use cluster from flag -C or --cluster
		// ============================

		cn := cmd.Flag("cluster").Value.String()
		if len(cn) != 0 {
			match := strings.Split(cn, ":")

			switch len(match) {
			case 1:
				cluster, err := storage.GetLocalCluster(cn)
				if err != nil {
					panic(err)
				}
				if cluster == nil {
					fmt.Println("cluster not found")
					return
				}
				host = cluster.Endpoint
			case 2:
				config.Headers = make(map[string]string, 0)
				config.Headers["X-Cluster-Name"] = cn
			default:
				fmt.Println("invalid data")
				return

			}

			cli.Cluster = client.NewClusterClient(host, config)
			ctx.SetClient(cli)
			return
		}

		// ============================
		// Use selected cluster
		// ============================

		cluster, err := storage.GetCluster()
		if err != nil {
			panic(err)
		}

		if cluster != "" {
			switch cluster[:2] {
			case "l.":
				cluster, err := storage.GetLocalCluster(cluster[2:])
				if err != nil {
					log.Error(err.Error())
					return
				}
				host = cluster.Endpoint
			case "r.":
				config.Headers = make(map[string]string, 0)
				config.Headers["X-Cluster-Name"] = cluster[2:]
			default:
				fmt.Println("can not read cluster info: check cache data ($HOME/.lastbackend)")
			}
		}

		cli.Cluster = client.NewClusterClient(host, config)
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
			if err := cmd.Help(); err != nil {
				log.Error(err.Error())
				return
			}
			return
		}

		// Attach sub command for namespace
		ns.AddCommand(
			serviceCmd,
			secretCmd,
			routeCmd,
		)

		if err := ns.Execute(); err != nil {
			log.Error(err.Error())
			return
		}

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
		if _, err := fmt.Scan(&login); err != nil {
			log.Error(err.Error())
			return
		}

		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			log.Error(err.Error())
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
			fmt.Println(err.Error())
			return
		}

		if err := storage.SetToken(session.Token); err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Authorization successful!")
		return
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
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage your secret",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage your configs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Manage your volumes",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage your route",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage your cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage cluster nodes",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "Manage cluster ingress servers",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Manage cluster discovery servers",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

// Execute adds all child commands to the root command
func Execute() {

	cobra.OnInitialize()

	var getSSLPath = func(filepath string) string {
		return strings.Join([]string{filesystem.HomeDir(), ".lastbackend", filepath}, string(os.PathSeparator))
	}

	RootCmd.PersistentFlags().StringP("cluster", "C", "", "Use cluster for operations")
	RootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	RootCmd.PersistentFlags().Bool("tls", false, "Use TLS")
	RootCmd.PersistentFlags().String("tlscacert", getSSLPath("ca.pem"), fmt.Sprintf("Trust certs signed only by this CA (default \"%s\")", getSSLPath("ca.pem")))
	RootCmd.PersistentFlags().String("tlscert", getSSLPath("cert.pem"), fmt.Sprintf("Path to TLS certificate file (default \"%s\")", getSSLPath("cert.pem")))
	RootCmd.PersistentFlags().String("tlskey", getSSLPath("key.pem"), fmt.Sprintf("Path to TLS key file (default \"%s\")", getSSLPath("key.pem")))
	RootCmd.PersistentFlags().String("token", "", "Set access token for header")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
