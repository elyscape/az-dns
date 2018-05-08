package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/elyscape/az-dns/helpers"
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
    az-dns clear A example.com -z example.com
        Removes the A record at the apex of example.com
    az-dns clear TXT sub -z example.com
        Removes the TXT record for sub.example.com
    az-dns clear NS sub.example.com.example.com -r -z example.com
        Removes the NS record for sub.example.com.example.com`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		recordType := dns.RecordType(strings.ToUpper(args[0]))
		hostname := args[1]

		client, err := helpers.NewRecordSetClient(dns.DefaultBaseURI)
		if err != nil {
			return
		}

		resourceGroup := viper.GetString("resource-group")
		if resourceGroup == "" {
			err = fmt.Errorf("a resource group name is required")
			return
		}

		zone := viper.GetString("zone")
		if zone == "" {
			err = fmt.Errorf("a DNS zone name is required")
			return
		}

		cmd.SilenceUsage = true

		relative := viper.GetBool("relative")
		recordName := helpers.GenerateRecordName(hostname, zone, relative)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		_, err = client.Delete(ctx, resourceGroup, zone, recordName, recordType, "")

		if err != nil {
			return
		}

		fmt.Println("success")

		return
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)

	clearCmd.PersistentFlags().BoolP("relative", "r", false, "HOSTNAME is a zone-relative label")
	if err := viper.BindPFlags(clearCmd.PersistentFlags()); err != nil {
		// This shouldn't happen
		panic(err)
	}
}
