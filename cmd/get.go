package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/elyscape/az-dns/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get TYPE HOSTNAME",
	Short: "Retrieve a DNS record set",
	Long: `Retrieve a record set from Azure DNS

This will print the contents of a particular record set on Azure DNS. The
currently-supported record types are A, AAAA, CAA, CNAME, and TXT. HOSTNAME may
be a fully-qualified domain name contained within the zone, a record name
relative to the zone, or either the empty string or @ for the apex. If a record
name contains the zone name (e.g. example.com.example.com), you should either
provide the FQDN or use the --relative flag.

Examples:
    az-dns get A example.com -z example.com
        Prints A records for example.com
    az-dns get AAAA sub -z example.com
        Prints AAAA records for sub.example.com
    az-dns get CNAME sub.example.com -r -z example.com
        Prints the CNAME record for sub.example.com.example.com`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		recordType := dns.RecordType(strings.ToUpper(args[0]))
		hostname := args[1]

		client, err := helpers.NewRecordSetClient(dns.DefaultBaseURI)
		if err != nil {
			return err
		}

		resourceGroup := viper.GetString("resource-group")
		if resourceGroup == "" {
			return fmt.Errorf("a resource group name is required")
		}

		zone := viper.GetString("zone")
		if zone == "" {
			return fmt.Errorf("a DNS zone name is required")
		}

		cmd.SilenceUsage = true

		relative := viper.GetBool("relative")
		recordName := helpers.GenerateRecordName(hostname, zone, relative)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		rrset, err := client.Get(ctx, resourceGroup, zone, recordName, recordType)

		if err != nil {
			return err
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
		case dns.CAA:
			if rrset.CaaRecords != nil {
				for _, record := range *rrset.CaaRecords {
					fmt.Printf("%v %v %q\n", *record.Flags, *record.Tag, *record.Value)
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
			out, err := json.Marshal(rrset.RecordSetProperties)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", out)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().BoolP("relative", "r", false, "HOSTNAME is a zone-relative label")
	if err := viper.BindPFlags(getCmd.PersistentFlags()); err != nil {
		// This shouldn't happen
		panic(err)
	}
}
