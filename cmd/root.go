package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ggit",
	Short: "Speed up the repo cloning from the github.com",
	Run: func(cmd *cobra.Command, args []string) {
		versionBool, err := cmd.Flags().GetBool("version")
		if err != nil {
			os.Exit(1)
		}
		if versionBool {
			fmt.Println(version)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		cloneCmd,
		versionCmd,
	)

	rootCmd.Flags().BoolP("version", "v", false, "Prints the version of ggit")
}
