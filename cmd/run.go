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

		generate, err := cmd.Flags().GetBool("generate")
		if err != nil {
			pterm.Error.Println(err)
			panic(err)
		}

		if generate {
			pterm.Info.Println("Generating config file")
			cyclopes.GeneratorCLI()
		} else {
			if cfgFile != "" {
				pterm.Info.Println("Using config file:", cfgFile)
			} else {
				pterm.Info.Println("Using default config file:", "./cyclops.yml")
				cfgFile = "./cyclops.yml"
			}
			pterm.FgLightYellow.Println("🧿 Starting visual testing")
			cyclopes.Start(cfgFile)
		}

	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ${pwd}/cyclops.yml)")
	rootCmd.PersistentFlags().BoolP("generate", "g", false, "generate a config file with more zing 🚀")
}

func initConfig() {

}
