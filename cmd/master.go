package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/michaelhenkel/ckube/kuberesources"
	"github.com/michaelhenkel/ckube/run"
	"github.com/michaelhenkel/ckube/utils"
	"github.com/michaelhenkel/remoteExec/client/executor"
	"github.com/spf13/cobra"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

var MasterMemory int
var MasterNetwork string
var MasterCpus int
var MasterDisk string
var CreateContrail bool

func init() {
	masterAddCmd.Flags().IntVarP(&MasterMemory, "memory", "m", 8000, "memory")
	masterAddCmd.Flags().IntVarP(&MasterCpus, "cpus", "c", 4, "cpus")
	masterAddCmd.Flags().StringVarP(&MasterDisk, "disk", "d", "15G", "disk")
	masterAddCmd.Flags().StringVarP(&MasterNetwork, "net", "v", "vpnkit", "network mode")
	masterAddCmd.Flags().BoolVarP(&CreateContrail, "contrail", "x", false, "create contrail")
	addCmd.AddCommand(masterAddCmd)
}

type Master struct {
	ClusterName    string
	InternalIP     string
	ExternalIP     string
	Name           string
	Cpus           int
	Memory         int
	Disk           string
	CreateContrail bool
	Network        string
}

var masterAddCmd = &cobra.Command{
	Use:   "master",
	Short: "master",
	Long:  `master`,
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		aCmd := cmd.Parent()
		clusterName := aCmd.Parent().Name()

		master := &Master{
			ClusterName:    clusterName,
			Name:           args[0],
			Cpus:           MasterCpus,
			Memory:         MasterMemory,
			Disk:           MasterDisk,
			CreateContrail: CreateContrail,
			Network:        MasterNetwork,
		}

		if err := master.addMaster(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		masterPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", master.ClusterName, "master", master.Name)
		socketPath := masterPath + "/00000003.00000947"

		err := utils.Retry(10, 2*time.Second, func() (err error) {
			fmt.Println("Waiting for socket path")
			_, err = os.Stat(socketPath)
			return
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		e := &executor.Executor{
			Socket: socketPath,
		}

		var ipResult *string
		fmt.Println("getting ip address from socket ")
		err = utils.Retry(40, 2*time.Second, func() (err error) {
			fmt.Println("Trying to get IP")
			ipResult, err = e.GetIP()
			return
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ip := []byte(*ipResult + "\n")
		err = ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", master.ClusterName, "master", master.Name, "vm.ip"), ip, 0644)
		master.InternalIP = *ipResult
		if master.Network == "vmnet" {
			master.ExternalIP = *ipResult
		}
		if master.Network == "vpnkit" {
			master.ExternalIP = "127.0.0.1"
		}

		if err := master.waitForK3S(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var tokenResult *string
		err = utils.Retry(10, 2*time.Second, func() (err error) {
			fmt.Println("Trying to get /var/lib/rancher/k3s/server/token")
			file := "/var/lib/rancher/k3s/server/token"
			tokenResult, err := e.GetFileContent(file)
			if *tokenResult == "" {
				err = fmt.Errorf("couldn't get content of /var/lib/rancher/k3s/server/token")
			}
			return
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		k3stoken := []byte(*tokenResult + "\n")
		err = ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", master.ClusterName, "master", master.Name, "node-token"), k3stoken, 0644)

		var k3sConfigResult *string
		err = utils.Retry(10, 2*time.Second, func() (err error) {
			fmt.Println("Trying to get /containers/services/k3s/rootfs/etc/rancher/k3s/k3s.yaml")
			file := "/containers/services/k3s/rootfs/etc/rancher/k3s/k3s.yaml"
			k3sConfigResult, err := e.GetFileContent(file)
			if *k3sConfigResult == "" {
				err = fmt.Errorf("couldn't get content of /containers/services/k3s/rootfs/etc/rancher/k3s/k3s.yaml")
			}
			return
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		k3sconfig := strings.Replace(string(*k3sConfigResult), "127.0.0.1", string(ip), -1)
		k3sconfigByte := []byte(k3sconfig)
		err = ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", master.ClusterName, "master", master.Name, "k3s.yaml"), k3sconfigByte, 0644)
		if err != nil {
			log.Fatalf("Cannot write k3s config to file: %v", err)
		}

		fmt.Println("export KUBECONFIG=" + masterPath + "/k3s.yaml")

		if !master.CreateContrail {
			if err = kuberesources.CreateContrailResources(masterPath + "/k3s.yaml"); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	},
}

func (m *Master) addMaster() error {
	var args []string
	masterPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", m.ClusterName, "master", m.Name)
	kernelpath := filepath.Join(os.Getenv("HOME"), ".ckube", "images")

	if _, err := os.Stat(masterPath); os.IsNotExist(err) {
		err := os.MkdirAll(masterPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create %q\n", masterPath)
			return err
		}
		fmt.Printf("Created Node in %q\n", masterPath)
	} else {
		fmt.Printf("Node %s already exists in %q\n", m.Name, masterPath)
	}

	cliArgs := `{
			"cliargs": {
			  "entries": {
				"args": {
				  "content": "server --cluster-cidr=10.32.0.0/12 --service-cidr=10.96.0.0/12 --no-flannel --no-deploy=traefik"
				}
			  }
			}
		}`
	cliArgsPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", m.ClusterName, "master", m.Name, "cliargs.json")
	err := ioutil.WriteFile(cliArgsPath, []byte(cliArgs), 0644)
	if err != nil {
		log.Fatalf("Cannot write cliArgs to file: %v", err)
	}
	args = append(args, "-cpus="+strconv.Itoa(m.Cpus))
	args = append(args, "-mem="+strconv.Itoa(m.Memory))
	args = append(args, "-disk=file="+masterPath+"/disk2.img,size="+m.Disk+",format=qcow2")
	args = append(args, "-console-file")
	args = append(args, "-kernelpath="+kernelpath)
	args = append(args, "-kernelprefix=ckube")
	args = append(args, "-state="+masterPath)
	args = append(args, "-data-file="+cliArgsPath)
	args = append(args, "-vsock-ports=2375")
	if m.Network == "vpnkit" {
		args = append(args, "-networking=vpnkit")
		//args = append(args, "-publish=6443:6443/tcp")
	}
	if m.Network == "vmnet" {
		args = append(args, "-networking=vmnet")
	}
	run.Run(args)

	return nil
}

func (m *Master) waitForK3S() error {

	err := utils.Retry(10, 2*time.Second, func() (err error) {
		fmt.Println("Waiting for k3s to come up on " + m.ExternalIP + ":6443")
		_, err = net.Dial("tcp", m.ExternalIP+":6443")

		return
	})
	if err != nil {
		return err
	}

	return nil
}
