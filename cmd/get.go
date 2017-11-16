package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/elyscape/azure-dns-client/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get TYPE HOSTNAME",
	Short: "Retrieve a DNS record set",
	Long: `Retrieve a record set from Azure DNS

This will print the contents of a particular record set on Azure DNS. The
currently-supported record types are A, AAAA, CNAME, and TXT. HOSTNAME may be a
fully-qualified domain name contained within the zone, a record name relative
to the zone, or either the empty string or @ for the apex. If a record name
contains the zone name (e.g. example.com.example.com), you should either
provide the FQDN or use the --relative flag.

Examples:
	azure-dns-client get A example.com -z example.com
		Prints A records for example.com
	azure-dns-client get AAAA sub -z example.com
		Prints AAAA records for sub.example.com
	azure-dns-client get CNAME sub.example.com -r -z example.com
		Prints the CNAME record for sub.example.com.example.com`,
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

		rrset, err := client.Get(resourceGroup, zone, recordName, recordType)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch recordType {
		case dns.A:
			if rrset.ARecords != nil {
				for _, record := range *rrset.ARecords {
					fmt.Println(*record.Ipv4Address)
				}
			}
		case dns.AAAA:
			if rrset.AaaaRecords != nil {
				for _, record := range *rrset.AaaaRecords {
					fmt.Println(*record.Ipv6Address)
				}
			}
		case dns.CNAME:
			if rrset.CnameRecord != nil {
				record := *rrset.CnameRecord
				fmt.Println(*record.Cname)
			}
		case dns.TXT:
			if rrset.TxtRecords != nil {
				for _, record := range *rrset.TxtRecords {
					for _, line := range *record.Value {
						fmt.Println(line)
					}
				}
			}
		default:
			if out, err := json.Marshal(rrset.RecordSetProperties); err == nil {
				fmt.Printf("%s\n", out)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().BoolP("relative", "r", false, "HOSTNAME is a zone-relative label")
	viper.BindPFlags(getCmd.PersistentFlags())
}