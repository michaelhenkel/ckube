package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/michaelhenkel/ckube/utils"
	"github.com/michaelhenkel/remoteExec/client/executor"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(keyGetCmd)
}

var keyGetCmd = &cobra.Command{
	Use:   "key",
	Short: "key",
	Long:  `key`,
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		//aCmd := cmd.Parent()
		clusterName := cmd.Parent().Parent().Name()
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
		var pubKey *string
		fmt.Println("getting pub key from socket ")
		err = utils.Retry(40, 2*time.Second, func() (err error) {
			fmt.Println("Trying to get pub key")
			pubKey, err = e.ExecuteCommand("/bin/cat /id_rsa.pub")
			return
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(*pubKey)
	},
}
