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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"wumpgo.dev/wumpgo/wumpgoctl/internal/scaffolding"
)

type tmplArgs struct {
	Package string
	HTTP    bool
	Gateway bool
	BotName string
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new bot built around wumpgo",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := scaffolding.Templates()

		pkg, _ := cmd.Flags().GetString("pkg")
		root, _ := cmd.Flags().GetString("root")
		http, _ := cmd.Flags().GetBool("http")
		gateway, _ := cmd.Flags().GetBool("gateway")
		name, _ := cmd.Flags().GetString("name")

		tArgs := tmplArgs{
			Package: pkg,
			HTTP:    http,
			Gateway: gateway,
			BotName: name,
		}

		for filename, tmpl := range t {
			filename = root + "/" + filename

			buf := &bytes.Buffer{}
			err := tmpl.Execute(buf, tArgs)
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

				f.Write(buf.Bytes())

				f.Close()
			}
		}

		return nil
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
}
