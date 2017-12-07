package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/elyscape/azure-dns-client/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear TYPE HOSTNAME",
	Short: "Delete a DNS record set",
	Long: `Delete a record set from Azure DNS

This will remove a record set from Azure DNS. HOSTNAME may be a fully-qualified
domain name contained within the zone, a record name relative to the zone, or
either the empty string or @ for the apex. If a record name contains the zone
name (e.g. example.com.example.com), you should either provide the FQDN or use
the --relative flag.

Examples:
    azure-dns-client clear A example.com -z example.com
        Removes the A record at the apex of example.com
    azure-dns-client clear TXT sub -z example.com
        Removes the TXT record for sub.example.com
    azure-dns-client clear NS sub.example.com.example.com -r -z example.com
        Removes the NS record for sub.example.com.example.com`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		recordType := dns.RecordType(strings.ToUpper(args[0]))
		hostname := args[1]

		client, err := helpers.NewRecordSetClient(dns.DefaultBaseURI)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		resourceGroup := viper.GetString("resource-group")
		if resourceGroup == "" {
			fmt.Println("A resource group name is required.")
			os.Exit(1)
		}

		zone := viper.GetString("zone")
		if zone == "" {
			fmt.Println("A DNS zone name is required.")
			os.Exit(1)
		}

		relative := viper.GetBool("relative")
		recordName := helpers.GenerateRecordName(hostname, zone, relative)

		_, err = client.Delete(resourceGroup, zone, recordName, recordType, "")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Success.")
	},
}

func init() {
	RootCmd.AddCommand(clearCmd)

	clearCmd.PersistentFlags().BoolP("relative", "r", false, "HOSTNAME is a zone-relative label")
	viper.BindPFlags(clearCmd.PersistentFlags())
}
