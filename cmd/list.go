package cmd

import (
	"fmt"

	"github.com/andrei-don/multi-k8s/k8s"
	"github.com/andrei-don/multi-k8s/multipass"
	"github.com/spf13/cobra"
)

func listCommand(multipassListFunction func() (string, error), k8sFilterNodesListCmdFunction func(string) string) {
	CheckMultipass()
	multipassList, err := multipassListFunction()
	if err != nil {
		fmt.Println("Error listing multipass nodes:", err)
		return
	}
	if k8sFilterNodesListCmdFunction(multipassList) == "" {
		fmt.Println("You have no multi-k8s clusters!")
	} else {
		fmt.Println(k8sFilterNodesListCmdFunction(multipassList))
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Command for listing cluster nodes",
	Long:  `Use this command to list the nodes from your cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		listCommand(multipass.List, k8s.FilterNodesListCmd)
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
