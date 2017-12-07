package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "azure-dns-client",
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

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.azure-dns-client.yaml)")

	// credentials
	RootCmd.PersistentFlags().String("client-id", "", "Azure client ID")
	RootCmd.PersistentFlags().String("client-secret", "", "Azure client secret")
	RootCmd.PersistentFlags().String("tenant-id", "", "Azure tenant ID")
	RootCmd.PersistentFlags().String("subscription-id", "", "Azure subscription ID")

	// resource info
	RootCmd.PersistentFlags().StringP("resource-group", "g", "", "Name of the resource group")
	RootCmd.PersistentFlags().StringP("zone", "z", "", "Name of the DNS zone")

	// other
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	viper.BindPFlags(RootCmd.PersistentFlags())
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

		// Search config in home directory with name ".azure-dns-client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".azure-dns-client")
	}

	viper.SetEnvPrefix("AZURE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && viper.GetBool("verbose") {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
