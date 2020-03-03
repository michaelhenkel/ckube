package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/michaelhenkel/ckube/run"
	"github.com/spf13/cobra"
)

func init() {
	nodeAddCmd.Flags().IntVarP(&Memory, "memory", "m", 4000, "memory")
	nodeAddCmd.Flags().IntVarP(&Cpus, "cpus", "c", 2, "cpus")
	nodeAddCmd.Flags().StringVarP(&Disk, "disk", "d", "15G", "disk")
	addCmd.AddCommand(nodeAddCmd)
}

type Node struct {
	ClusterName string
	IP          string
	Name        string
	Cpus        int
	Memory      int
	Disk        string
}

var nodeAddCmd = &cobra.Command{
	Use:   "node",
	Short: "node",
	Long:  `node`,
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		aCmd := cmd.Parent()
		clusterName := aCmd.Parent().Name()

		node := &Node{
			ClusterName: clusterName,
			Name:        args[0],
			Cpus:        Cpus,
			Memory:      Memory,
			Disk:        Disk,
		}

		if err := node.addNode(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		/*
			if err := kuberesources.CreateContrailResources(masterPath + "/k3s.yaml"); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		*/

	},
}

func (m *Node) addNode() error {
	var args []string
	mastersPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", m.ClusterName, "master")
	masters, err := ioutil.ReadDir(mastersPath)
	master := masters[0].Name()
	masterPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", m.ClusterName, "master", master)
	masterIP, err := ioutil.ReadFile(masterPath + "/vm.ip")
	masterIPTrimmed := strings.TrimSuffix(string(masterIP), "\n")
	masterToken, err := ioutil.ReadFile(masterPath + "/node-token")
	masterTokenTrimmed := strings.TrimSuffix(string(masterToken), "\n\n")
	nodePath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", m.ClusterName, "nodes", m.Name)
	kernelpath := filepath.Join(os.Getenv("HOME"), ".ckube", "images")

	if _, err := os.Stat(nodePath); os.IsNotExist(err) {
		err := os.MkdirAll(nodePath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create %q\n", nodePath)
			return err
		}
		fmt.Printf("Created Node in %q\n", nodePath)
	} else {
		fmt.Printf("Node %s already exists in %q\n", m.Name, nodePath)
	}

	cliArgs := "{\n" +
		"\"cliargs\": {\n" +
		"\"entries\": {\n" +
		"\"args\": {\n" +
		"\"content\": \"agent --server https://" + masterIPTrimmed + ":6443 --token " + masterTokenTrimmed + "\"\n" +
		"}\n" +
		"}\n" +
		"}\n" +
		"}\n"
	cliArgsPath := filepath.Join(nodePath, "cliargs.json")
	err = ioutil.WriteFile(cliArgsPath, []byte(cliArgs), 0644)
	if err != nil {
		log.Fatalf("Cannot write cliArgs to file: %v", err)
	}
	args = append(args, "-cpus="+strconv.Itoa(m.Cpus))
	args = append(args, "-mem="+strconv.Itoa(m.Memory))
	args = append(args, "-disk=file="+nodePath+"/disk2.img,size="+m.Disk+",format=qcow2")
	args = append(args, "-networking=vmnet")
	args = append(args, "-console-file")
	args = append(args, "-kernelpath="+kernelpath)
	args = append(args, "-kernelprefix=ckube")
	args = append(args, "-state="+nodePath)
	args = append(args, "-data-file="+cliArgsPath)
	run.Run(args)
	ipAddress, err := ioutil.ReadFile(nodePath + "/vm.ip")
	if err != nil {
		return err
	}
	ipAddressTrimmed := strings.TrimSuffix(string(ipAddress), "\n")
	m.IP = string(ipAddressTrimmed)
	return nil
}
