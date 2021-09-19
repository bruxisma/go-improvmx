package improvmx_test

import (
	"context"
	"fmt"
	"log"

	"occult.work/improvmx"
)

func Example_account() {
	session, error := improvmx.New("token")
	if error != nil {
		log.Fatal(error)
	}
	account, error := session.Account.Read(context.Background())
	if error != nil {
		log.Fatal(error)
	}
	fmt.Printf("Name: %s\tPlan: %s\n", account.CompanyName, account.Plan.Name)
}

func Example_aliases() {
	session, error := improvmx.New("token")
	if error != nil {
		log.Fatal(error)
	}
	options := improvmx.NewListOption().
		SetIsActive(true).
		SetStartsWith("r")
	aliases, error := session.Aliases.List(context.Background(), "example.com", options)
	if error != nil {
		log.Fatal(error)
	}
	for _, alias := range aliases {
		fmt.Printf("%s <%s>\n", alias.Name, alias.Address)
	}
}

func Example_credentials() {
	session, error := improvmx.New("token")
	if error != nil {
		log.Fatal(error)
	}
	ctx := context.Background()
	account, error := session.Account.Read(ctx)
	if error != nil {
		log.Fatal(error)
	}
	if !account.Premium {
		log.Fatal("SMTP Credentials are a premium account feature")
	}
	credentials, error := session.Credentials.List(ctx, "example.com")
	if error != nil {
		log.Fatal(error)
	}
	for _, credential := range credentials {
		fmt.Println(credential.Username)
	}
}

func Example_domains() {

}

func Example_logs() {
	session, error := improvmx.New("token")
	if error != nil {
		log.Fatal(error)
	}
	logs, error := session.Aliases.Logs(context.Background(), "example.com", "richard")
	if error != nil {
		log.Fatal(error)
	}

	for _, entry := range logs {
		fmt.Printf("ID: %s, Hostname: %s, Subject: %s",
			entry.ID,
			entry.Hostname,
			entry.Subject)
	}
}
