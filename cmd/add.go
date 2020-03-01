package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {

}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long:  `add`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Adds " + strings.Join(args, " "))
	},
}
