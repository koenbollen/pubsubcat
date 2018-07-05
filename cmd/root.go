package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string
var globalProjectID string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pubsubcat",
	Short: "The Google Pub/Sub Swiss Army Knife",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auto detecting publish/subscribe (not supported yet)")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.SetFlags(0)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.pubsubcat)")
	rootCmd.PersistentFlags().StringVarP(&globalProjectID, "project", "p", "", "Google Cloud Project to work under")

	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "only output messages")
	rootCmd.PersistentFlags().CountP("verbose", "v", "increase verbosity")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(path.Join(home, ".config"))
		viper.SetConfigName(".pubsubcat")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("] using config file:", viper.ConfigFileUsed())
	}
}

// GetVerbosity determines the level of verboseness based on the given
// flags. Quiet is 0 and default is 1, all -v flags add one.
func GetVerbosity(flags *pflag.FlagSet) int {
	if quiet, err := flags.GetBool("quiet"); err == nil && quiet {
		return 0
	}
	verbosity, _ := flags.GetCount("verbose")
	return verbosity + 1
}
