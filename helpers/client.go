package helpers

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	authfile "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/viper"
)

func NewRecordSetClient(baseURI string) (client dns.RecordSetsClient, err error) {
	var authorizer *autorest.BearerAuthorizer
	subscriptionID := viper.GetString("subscription-id")

	clientSetup, err := authfile.GetClientSetup(baseURI)
	if err == nil {
		authorizer = clientSetup.BearerAuthorizer
		subscriptionID = clientSetup.SubscriptionID
	} else {
		authorizer, err = GetAuthorizer(baseURI)
		if err != nil {
			return
		}
	}

	client = dns.NewRecordSetsClientWithBaseURI(baseURI, subscriptionID)
	client.Authorizer = authorizer

	return
}

func GetAuthorizer(baseURI string) (authorizer *autorest.BearerAuthorizer, err error) {
	credFields := []string{"client-id", "client-secret", "subscription-id", "tenant-id"}

	for _, field := range credFields {
		if !viper.IsSet(field) || viper.GetString(field) == "" {
			err = fmt.Errorf("Required credential option %v not provided.", field)
			return
		}
	}

	config, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, viper.GetString("tenant-id"))
	if err != nil {
		return
	}

	token, err := adal.NewServicePrincipalToken(*config, viper.GetString("client-id"), viper.GetString("client-secret"), baseURI)
	if err != nil {
		return
	}

	authorizer = autorest.NewBearerAuthorizer(token)
	return
}
