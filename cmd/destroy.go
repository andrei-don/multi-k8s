/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/andrei-don/multi-k8s/k8s"
	"github.com/andrei-don/multi-k8s/multipass"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Command for destroying a cluster",
	Long:  `Run this command to destroy your multi-k8s cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		multipassList, err := multipass.List()
		if err != nil {
			fmt.Println("Error listing multipass nodes:", err)
			return
		}
		k8s.DeleteClusterVMs(k8s.GetCurrentNodes(k8s.FilterNodesListCmd(multipassList)))
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
