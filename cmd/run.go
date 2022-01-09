package cmd

import (
	"github.com/konsalex/cyclopes/cmd/cyclopes"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "cyclops",
	Short: "Quick and dirty Visual Testing",
	Long: `
The _quick and dirty_ testing manifesto:
  1. Stop lying about your E2E testing, you do not have visual tests.
  2. Reduce the headaches of broken styles, just run visual tests.
  3. Be more confident that you ship good web products, again run visual tests.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.FgLightYellow.Println("🧿 Starting visual testing")
		cyclopes.Start(cfgFile)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ${pwd}/cyclops.yml)")
}

func initConfig() {
	if cfgFile != "" {
		pterm.Info.Println("Using config file:", cfgFile)
	} else {
		pterm.Info.Println("Using default config file:", "./cyclops.yml")
		cfgFile = "./cyclops.yml"
	}
}
