package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/michaelhenkel/remoteExec/client/executor"
	"github.com/spf13/cobra"
)

var HostPort int
var VMPort int
var Username string
var Address string

func init() {
	tunnelAddCmd.Flags().IntVarP(&HostPort, "hostport", "p", 6443, "host port")
	tunnelAddCmd.Flags().IntVarP(&VMPort, "vmport", "v", 6443, "vm port")
	tunnelAddCmd.Flags().StringVarP(&Username, "user", "u", "", "username")
	tunnelAddCmd.Flags().StringVarP(&Address, "address", "a", "192.168.64.1:22", "address")
	addCmd.AddCommand(tunnelAddCmd)
	tunnelDeleteCmd.Flags().IntVarP(&HostPort, "hostport", "p", 6443, "host port")
	tunnelDeleteCmd.Flags().IntVarP(&VMPort, "vmport", "v", 6443, "vm port")
	tunnelDeleteCmd.Flags().StringVarP(&Username, "user", "u", "", "username")
	tunnelDeleteCmd.Flags().StringVarP(&Address, "address", "a", "192.168.64.1:22", "address")
	deleteCmd.AddCommand(tunnelDeleteCmd)
}

var tunnelDeleteCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "tunnel",
	Long:  `tunnel`,
	Run: func(cmd *cobra.Command, args []string) {
		aCmd := cmd.Parent()
		clusterName := aCmd.Parent().Name()
		mastersPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", clusterName, "master")
		masters, err := ioutil.ReadDir(mastersPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		master := masters[0].Name()
		socketPath := mastersPath + "/" + master + "/00000003.00000947"
		e := &executor.Executor{
			Socket: socketPath,
		}
		result, err := e.DeleteTunnel(VMPort, HostPort, Username, Address)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(result)

	},
}

var tunnelAddCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "tunnel",
	Long:  `tunnel`,
	Run: func(cmd *cobra.Command, args []string) {
		aCmd := cmd.Parent()
		clusterName := aCmd.Parent().Name()
		mastersPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", clusterName, "master")
		masters, err := ioutil.ReadDir(mastersPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		master := masters[0].Name()
		socketPath := mastersPath + "/" + master + "/00000003.00000947"
		e := &executor.Executor{
			Socket: socketPath,
		}

		result, err := e.SetupTunnel(VMPort, HostPort, Username, Address)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(result)
	},
}
