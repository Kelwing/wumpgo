/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"wumpgo.dev/wumpgo/wumpgoctl/internal/cmdgen"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Parse comments on slash commands to generate interface implementations",
	Long: `Parses specially formatted comments on slash command structs
to generate interface implementations.
Example:
// MyCommand godoc
// @Name testcommand
// @Description Test base command
// @Name.en-US testcommand
// @Name.es-MX commandodepreueba
// @Type ChatInput
// @Permissions ManageRoles, KickMembers
type MyCommand struct {}
`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		logger := log.Level(zerolog.WarnLevel)
		ctx := logger.WithContext(context.Background())

		cmdgen.Gen(ctx, dir)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringP("dir", "d", ".", "Directory to process, you can omit this if your go:generate comment is in the same directory as the commands you want to process")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
