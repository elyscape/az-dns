package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cfgFile is the path to an optional configuration file.
var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "az-dns",
	Short: "Azure DNS record set manipulator",
	Long: `A simple command-line tool for manipulating Azure DNS record sets

This client provides an easy way to view and manipulate record sets in Azure
DNS. It authenticates as an Azure Active Directory service principal, using
credentials provided via:
    a. command-line flags
    b. environment variables
    c. a config file, or
    d. an Azure CLI auth file, with path specified in $AZURE_AUTH_LOCATION`,
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.az-dns.yaml)")

	// credentials
	rootCmd.PersistentFlags().String("client-id", "", "Azure client ID")
	rootCmd.PersistentFlags().String("client-secret", "", "Azure client secret")
	rootCmd.PersistentFlags().String("tenant-id", "", "Azure tenant ID")
	rootCmd.PersistentFlags().String("subscription-id", "", "Azure subscription ID")

	// resource info
	rootCmd.PersistentFlags().StringP("resource-group", "g", "", "Name of the resource group")
	rootCmd.PersistentFlags().StringP("zone", "z", "", "Name of the DNS zone")

	// other
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".az-dns" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".az-dns")
	}

	viper.SetEnvPrefix("AZURE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && viper.GetBool("verbose") {
		fmt.Println("using config file:", viper.ConfigFileUsed())
	}
}
