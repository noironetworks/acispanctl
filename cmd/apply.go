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
	"github.com/spf13/viper"
	"os"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the sessions deployment file",
	Long: `For applying span definition file related to ACI ERSPAN

E.g. acispanctl apply -f <filename>`,

	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("apply called")
		dfilename, _:= cmd.Flags().GetString("file")
		if dfilename == "" {
			fmt.Println("error: must specify -f")
			os.Exit(1)
		}
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.SetConfigFile(dfilename)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil { // Handle errors reading the config file
			//panic(fmt.Errorf("Fatal error config file: %s \n", err))
			fmt.Printf("error while opening the span definition file. %+v\n", err)
			os.Exit(1)
		}

		var C acispanctl.SpanConfig
		err = viper.Unmarshal(&C)
		if err != nil {
			fmt.Printf("unable to decode into struct, %+v\n", err)
			os.Exit(1)
		}
		err = acispanctl.ApplyVSPANConfig(C)
		if err != nil {
			fmt.Printf("unable to apply vspan config %+v\n", C, err)
			os.Exit(1)
		}
		//viper.WriteConfigAs("test.yaml")
	},
}



func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("file", "f", "", "Span Definition Filename")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
