/*
Copyright © 2022 Kelwing <kelwing@kelnet.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"wumpgo.dev/wumpgo/wumpgoctl/internal/config"
	"wumpgo.dev/wumpgo/wumpgoctl/internal/scaffolding"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a new bot built using wumpgo",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := scaffolding.Templates()

		cfg := config.Config{}

		if err := viper.Unmarshal(&cfg); err != nil {
			return err
		}

		v, err := semver.NewVersion(strings.Replace(runtime.Version(), "go", "", 1))
		if err != nil {
			return err
		}

		cfg.GoVersion = fmt.Sprintf("%d.%d", v.Major(), v.Minor())
		for filename, tmpl := range t {
			filename = viper.GetString("root") + "/" + filename

			buf := &bytes.Buffer{}
			err := tmpl.Execute(buf, &cfg)
			if err != nil {
				return err
			}

			if buf.Len() > 0 {
				// Check if the directory exists
				dir := filepath.Dir(filename)
				if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
					err := os.MkdirAll(dir, os.ModePerm)
					if err != nil {
						return err
					}
				}

				f, err := os.Create(filename)
				if err != nil {
					return err
				}

				if _, err := f.Write(buf.Bytes()); err != nil {
					return err
				}

				f.Close()
			}
		}

		// Save the config
		if err := viper.SafeWriteConfig(); err != nil {
			return err
		}

		tidyCmd := exec.Command("go", "mod", "tidy")

		return tidyCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().StringP("pkg", "p", "github.com/example/examplebot", "Base package for project")
	initCmd.Flags().StringP("root", "d", ".", "Root directory to create the project in")
	initCmd.Flags().BoolP("http", "w", false, "Include HTTP interactions support")
	initCmd.Flags().BoolP("gateway", "g", false, "Include gateway support")
	initCmd.Flags().StringP("name", "n", "ExampleBot", "Name of your bot")
	initCmd.Flags().Bool("nats", false, "Use NATS to forward gateway events")
	initCmd.Flags().Bool("redis", false, "Use Redis to forward gateway events")
	initCmd.Flags().Bool("local", false, "Use the local dispatcher for events")
	initCmd.Flags().Bool("codegen", false, "Use codegen for slash commands")
	initCmd.Flags().StringP("summary", "s", "A bot made with wumpgo (wumpgo.dev)", "Summary of the bots functionalty")
	initCmd.Flags().String("descrption", "", "Long description of the bots functionalty")
	initCmd.MarkFlagsMutuallyExclusive("nats", "redis", "local")
	viper.BindPFlag("meta.package", initCmd.Flags().Lookup("pkg"))
	viper.BindPFlag("meta.name", initCmd.Flags().Lookup("name"))
	viper.BindPFlag("features.http.enabled", initCmd.Flags().Lookup("http"))
	viper.BindPFlag("features.gateway.enabled", initCmd.Flags().Lookup("gateway"))
	viper.BindPFlag("features.gateway.nats", initCmd.Flags().Lookup("nats"))
	viper.BindPFlag("features.gateway.redis", initCmd.Flags().Lookup("redis"))
	viper.BindPFlag("features.gateway.local", initCmd.Flags().Lookup("local"))
	viper.BindPFlag("features.codegen", initCmd.Flags().Lookup("codegen"))
	viper.BindPFlag("root", initCmd.Flags().Lookup("root"))
	viper.BindPFlag("meta.summary", initCmd.Flags().Lookup("summary"))
	viper.BindPFlag("meta.description", initCmd.Flags().Lookup("descrption"))
}
