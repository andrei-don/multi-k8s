/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andrei-don/multi-k8s/multipass"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Command to list the nodes of your cluster.",
	Long:  `Use this command to list the nodes from your cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		multipassList, err := multipass.List()
		if err != nil {
			fmt.Println("Error listing multipass nodes:", err)
			return
		}
		fmt.Println(filterNodesList(multipassList))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Function below is used to filter through the 'multipass list' command and return only the cluster nodes (assuming that there are other unrelated multipass nodes as well)
func filterNodesList(multipassNodes string) string {
	lines := strings.Split(multipassNodes, "\n")

	var result []string

	re := regexp.MustCompile(`^(controller-node-[123]|worker-node-[123])\s+.*`)

	for _, line := range lines {
		if re.MatchString(line) || strings.HasPrefix(line, "Name") {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
