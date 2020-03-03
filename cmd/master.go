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
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var Memory int
var Cpus int
var Disk string

func init() {
	masterAddCmd.Flags().IntVarP(&Memory, "memory", "m", 8000, "memory")
	masterAddCmd.Flags().IntVarP(&Cpus, "cpus", "c", 4, "cpus")
	masterAddCmd.Flags().StringVarP(&Disk, "disk", "d", "15G", "disk")
	addCmd.AddCommand(masterAddCmd)
}

type Master struct {
	ClusterName string
	IP          string
	Name        string
	Cpus        int
	Memory      int
	Disk        string
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
			ClusterName: clusterName,
			Name:        args[0],
			Cpus:        Cpus,
			Memory:      Memory,
			Disk:        Disk,
		}

		if err := master.addMaster(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := master.waitForK3S(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err := utils.Retry(10, 2*time.Second, func() (err error) {
			fmt.Println("Waiting for k3s config to be written")
			err = master.getK3SConfig()
			return
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = utils.Retry(10, 2*time.Second, func() (err error) {
			fmt.Println("Waiting for k3s token to be written")
			err = master.getK3SToken()
			return
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		masterPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", master.ClusterName, "master", master.Name)
		fmt.Println("export KUBECONFIG=" + masterPath + "/k3s.yaml")

		if err = kuberesources.CreateContrailResources(masterPath + "/k3s.yaml"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func (m *Master) addMaster() error {
	var args []string
	masterPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", m.ClusterName, "master", m.Name)
	kernelpath := filepath.Join(os.Getenv("HOME"), ".ckube", "images")

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
	args = append(args, "-networking=vmnet")
	args = append(args, "-console-file")
	args = append(args, "-kernelpath="+kernelpath)
	args = append(args, "-kernelprefix=ckube")
	args = append(args, "-state="+masterPath)
	args = append(args, "-data-file="+cliArgsPath)
	run.Run(args)
	ipAddress, err := ioutil.ReadFile(masterPath + "/vm.ip")
	if err != nil {
		return err
	}
	ipAddressTrimmed := strings.TrimSuffix(string(ipAddress), "\n")
	m.IP = string(ipAddressTrimmed)
	return nil
}

func (m *Master) waitForK3S() error {

	err := utils.Retry(10, 2*time.Second, func() (err error) {
		fmt.Println("Waiting for k3s to come up on " + m.IP + ":6443")
		_, err = net.Dial("tcp", m.IP+":6443")

		return
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Master) getK3SConfig() error {
	agent, err := getAgent()
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", m.IP+":22", &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // FIXME: please be more secure in checking host keys
	})
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	results := make(chan string, 10)
	timeout := time.After(10 * time.Second)
	go func(session *ssh.Session) {
		results <- executeCmd("cat /containers/services/k3s/rootfs/etc/rancher/k3s/k3s.yaml", session)
	}(session)
	var k3sconfig []byte
	select {
	case res := <-results:
		res2 := strings.Replace(res, "127.0.0.1", m.IP, -1)
		k3sconfig = []byte(res2 + "\n")
		if len(strings.TrimSpace(string(k3sconfig))) == 0 {
			return fmt.Errorf("k3s config empty")
		}
	case <-timeout:
		return fmt.Errorf("timed out")
	}
	err = ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", m.ClusterName, "master", m.Name, "k3s.yaml"), k3sconfig, 0644)
	if err != nil {
		log.Fatalf("Cannot write k3s config to file: %v", err)
	}

	return nil
}

func (m *Master) getK3SToken() error {
	agent, err := getAgent()
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", m.IP+":22", &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // FIXME: please be more secure in checking host keys
	})
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	results := make(chan string, 10)
	timeout := time.After(10 * time.Second)
	go func(session *ssh.Session) {
		results <- executeCmd("ctr --namespace services.linuxkit tasks exec --exec-id ssh-xx k3s /bin/cat /var/lib/rancher/k3s/server/node-token", session)
	}(session)
	var k3stoken []byte
	select {
	case res := <-results:
		k3stoken = []byte(res + "\n")
		if len(strings.TrimSpace(string(k3stoken))) == 0 {
			return fmt.Errorf("k3s token empty")
		}
	case <-timeout:
		return fmt.Errorf("timed out")
	}
	err = ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", m.ClusterName, "master", m.Name, "node-token"), k3stoken, 0644)
	if err != nil {
		log.Fatalf("Cannot write k3s token to file: %v", err)
	}

	return nil
}

func executeCmd(command string, session *ssh.Session) string {
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)

	return fmt.Sprintf("%s", stdoutBuf.String())
}

func getAgent() (agent.Agent, error) {
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	return agent.NewClient(agentConn), err
}
