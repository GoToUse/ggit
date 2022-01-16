/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "ggit version: [v0.1]"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of ggit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
