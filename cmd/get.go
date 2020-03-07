package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gets " + strings.Join(args, " "))
	},
}
