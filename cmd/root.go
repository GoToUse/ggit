package cmd

import (
	"fmt"
	"log"
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
		versionCmd, // 获取版本号方式一
	)

	// 获取版本号方式二
	rootCmd.Flags().BoolP("version", "v", false, "Prints the version of ggit")

	// Initialize the configuration data.
	err := setConfig()
	if err != nil {
		log.Fatalf("init.setConfig err: %v", err)
	}

	// Assign data from the configuration file to variable.
	DefaultMirrorUrlMap = mirrorUrlArr
	DefaultGitPath, DefaultGithubUrl, DefaultGithubSuffix = GitC.FilePath, GitC.Website, GitC.UrlSuffix
}

// setConfig 将配置信息写入 GitC / mirrorUrlArr 中
func setConfig() error {
	config, err := NewConfig()
	if err != nil {
		return err
	}

	err = config.ReadSection("Git", &GitC)
	if err != nil {
		return err
	}

	err = config.ReadSection("MirrorUrl", &mirrorUrlArr)
	if err != nil {
		return err
	}

	return nil
}
