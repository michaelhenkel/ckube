package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	clusterCmd.AddCommand(createCmd)
	clusterCmd.AddCommand(deleteCmd)
	if err := getClusters(clusterCmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(clusterCmd)
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster",
	Long:  `cluster`,
	Args:  cobra.RangeArgs(1, 3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cluster " + strings.Join(args, " "))
	},
}

func getClusters(clusterCmd *cobra.Command) error {
	clusterPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters")
	files, _ := ioutil.ReadDir(clusterPath)
	for _, f := range files {
		if f.IsDir() {
			var ccmd = &cobra.Command{
				Use:   f.Name(),
				Short: f.Name(),
				Long:  f.Name(),
				Args:  cobra.RangeArgs(1, 3),
				Run: func(cmd *cobra.Command, args []string) {
					fmt.Println("Cluster " + strings.Join(args, " "))
				},
			}

			ccmd.AddCommand(addCmd)
			ccmd.AddCommand(deleteCmd)
			ccmd.AddCommand(getCmd)
			clusterCmd.AddCommand(ccmd)
		}
	}
	return nil
}
