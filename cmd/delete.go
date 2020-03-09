package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/michaelhenkel/remoteExec/client/executor"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	Long:  `delete`,
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := cmd.Parent().Name()
		//clusterName := aCmd.Parent().Name()
		if err := deleteCluster(clusterName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func deleteCluster(clusterName string) error {
	mastersPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", clusterName, "master")
	masters, err := ioutil.ReadDir(mastersPath)
	if err != nil {
		return err
	}
	master := masters[0].Name()
	masterPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", clusterName, "master", master)
	socketPath := masterPath + "/00000003.00000947"
	masterPathFile, err := ioutil.ReadDir(masterPath)
	if err != nil {
		return err
	}
	e := executor.Executor{
		Socket: socketPath,
	}
	var vpnkitPid *int
	var hyperkitPid *int
	for _, file := range masterPathFile {
		if file.Name() == "vpnkit.pid" {
			pidByte, err := ioutil.ReadFile(masterPath + "/" + file.Name())
			if err != nil {
				return err
			}
			pidString := string(pidByte)
			_vpnkitPid, err := strconv.Atoi(strings.Trim(pidString, "\n"))
			vpnkitPid = &_vpnkitPid
			if err != nil {
				return err
			}
		}
		if file.Name() == "hyperkit.pid" {
			pidByte, err := ioutil.ReadFile(masterPath + "/" + file.Name())
			if err != nil {
				return err
			}
			pidString := string(pidByte)
			_hyperkitPid, err := strconv.Atoi(strings.Trim(pidString, "\n"))
			hyperkitPid = &_hyperkitPid
			if err != nil {
				return err
			}
		}
	}
	e.ExecuteCommand("/sbin/poweroff -f")
	if hyperkitPid != nil {
		proc, err := os.FindProcess(*hyperkitPid)
		if err == nil {
			if err := proc.Kill(); err != nil {
				return err
			}
		}

	}
	if vpnkitPid != nil {
		proc, err := os.FindProcess(*vpnkitPid)
		if err == nil {
			if err := proc.Kill(); err != nil {
				return err
			}
		}
	}
	err = os.RemoveAll(masterPath)
	if err != nil {
		return err
	}
	return nil
}
