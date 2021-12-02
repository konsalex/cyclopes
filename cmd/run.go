package cmd

import (
	"fmt"

	"github.com/konsalex/cyclopes/cmd/cyclopes"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "cyclops",
	Short: "Quick and dirty Visual Testing",
	Long: `
The _quick and dirty_ testing manifesto:
  1. Stop lying about your code, you do not write tests.
  2. Save some time from your day-to-day work, by at least run a visual test.
  3. Be more confident that you build good things, and eventually you will become a millionaire.	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧿 Starting visual testing")
		cyclopes.Start(cfgFile)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ${pwd}/cyclops.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		fmt.Println("Using config file:", cfgFile)
	} else {
		fmt.Println("Using default config file:", "./cyclops.yml")
		cfgFile = "./cyclops.yml"
	}
}
