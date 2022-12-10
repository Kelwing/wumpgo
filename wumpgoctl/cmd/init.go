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
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"wumpgo.dev/wumpgo/wumpgoctl/internal/scaffolding"
)

type tmplArgs struct {
	Package string
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new bot built around wumpgo",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := scaffolding.Templates()

		pkg, _ := cmd.Flags().GetString("pkg")
		root, _ := cmd.Flags().GetString("root")

		tArgs := tmplArgs{
			Package: pkg,
		}

		for filename, tmpl := range t {
			filename = root + "/" + filename
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

			_ = tmpl.Execute(f, tArgs)

			f.Close()
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
}
