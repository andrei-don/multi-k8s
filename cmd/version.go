/*
Copyright © 2024 Alex Stan
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Command for showing the version for multi-k8s",
	Long:  `Use this command to get the version for multi-k8s.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("multi-k8s version: %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
