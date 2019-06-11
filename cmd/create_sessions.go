// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"acispanctl/pkg/acispanctl"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// sessionsCmd represents the sessions command
var createSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		cont, _:= cmd.Flags().GetString("container")
		if cont == "" {
			fmt.Println("specify name of a container to span using -c")
			return
		}

		dest, _:= cmd.Flags().GetString("destination")
		if dest == "" {
			fmt.Println("specify ER SPAN destination using -t")
			return
		}

		domain, _:= cmd.Flags().GetString("domain")
		if domain == "" {
			fmt.Println("WARN: using default domain kube")
			domain = "kube"
		}

		namespace, _:= cmd.Flags().GetString("namespace")
		if namespace == "" {
			fmt.Println("WARN: using default namespace default")
			namespace = "default"
		}


		err := acispanctl.CreateSpanSessionFromCont(cont, domain, namespace, dest)
		if err != nil { // Handle errors creating the config file
			fmt.Printf("error while creating the span definition file. %+v\n", err)
			os.Exit(1)
		}
		fname := fmt.Sprintf("%s-vspan.yaml", cont)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.SetConfigFile(fname)
		err = viper.ReadInConfig() // Find and read the config file
		if err != nil { // Handle errors reading the config file
			fmt.Printf("error while opening the span definition file. %+v\n", err)
			os.Exit(1)
		}

		var C acispanctl.SpanConfig
		err = viper.Unmarshal(&C)
		if err != nil {
			fmt.Printf("unable to decode into struct, %+v", err)
			os.Exit(1)
		}
		err = acispanctl.ApplyVSPANConfig(C)
		if err != nil {
			fmt.Printf("unable to apply vspan config %+v", C, err)
			os.Exit(1)
		}
	},
}

func init() {
	createCmd.AddCommand(createSessionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sessionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sessionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
