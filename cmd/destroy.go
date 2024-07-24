/*
Copyright © 2024 Alex Stan
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/andrei-don/multi-k8s/k8s"
	"github.com/andrei-don/multi-k8s/multipass"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Command for destroying a cluster",
	Long:  `Run this command to destroy your multi-k8s cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		multipassList, err := multipass.List()
		if err != nil {
			fmt.Println("Error listing multipass nodes:", err)
			return
		}
		nodesList := k8s.FilterNodesListCmd(multipassList)
		if nodesList == "" {
			fmt.Println("There are no cluster nodes!")
		}
		k8s.DeleteClusterVMs(k8s.GetCurrentNodes(nodesList))
		time.Sleep(2 * time.Second)
		fmt.Println()
		fmt.Println("Removing the contents of the dhcpd_leases file (requires sudo)... DO NOT PROCEED if you have other VMs running on your local apart from the kubernetes ones! (y/n)? ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		input = strings.TrimSpace(input)

		if input == "y" {
			cmd := exec.Command("sudo", "truncate", "-s", "0", "/var/db/dhcpd_leases")
			_, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
				return
			}
		} else {
			fmt.Println("dhcpd_leases file was not modified")
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
