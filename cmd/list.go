/*
Copyright © 2024 Alex Stan
*/
package cmd

import (
	"fmt"

	"github.com/andrei-don/multi-k8s/k8s"
	"github.com/andrei-don/multi-k8s/multipass"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Command for listing nodes of your cluster.",
	Long:  `Use this command to list the nodes from your cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		multipassList, err := multipass.List()
		if err != nil {
			fmt.Println("Error listing multipass nodes:", err)
			return
		}
		fmt.Println(k8s.FilterNodesListCmd(multipassList))
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
