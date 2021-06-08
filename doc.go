// Copyright 2021 Isabella Muerte. All rights reserved.
//
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE.md file that accompanied this package

/*
ImprovMX is a software as a service that allows setting up email forwarding in
seconds and comes with free, premium, and business tiers. This package is able
to wrap the entire API documented https://improvmx.com/api however not all
features will work across all tiers. At the present time, only the v3 API is
supported.

Within this documentation, the REST API is refered to the as the ImprovMX REST
API, and the golang API is simply called either the golang API or just "the
API".

Usage

The improvmx package is designed to be used with go modules enabled only.

	import "occult.work/improvmx"

Examples

Construct a new API session, then access different parts of the ImprovMX REST
API. For example:

	session, error := improvmx.New("authentication-token")
	// List all aliases for domain "example.com"
	aliases, error := session.Aliases.List(context.Background(), "example.com")
	for alias := range aliases {
		fmt.Printf("%s: %s\n", alias.Name, alias.Address)
	}

Endpoints

The ImprovMX REST API has several endpoints to access information regarding
Domains, Aliases, SMTP Credentials, and Logs for both Domains and Aliases.
Unlike the ImprovMX REST API, however, the golang API has organized these
endpoints into Domain, Aliases, Account, and Credentials. Every operation
within these endpoints attempts to follow the CRUD (Create, Read, Update,
Delete) naming convention. There are some outliers, such as
Account.Whitelabels, or Domain.Verify. Additionally, logs are grouped with
their specific endpoints. Users are also able to List some resources:

	domains, error := session.Domains.List()

Due to limitations with the ImprovMX REST API, the Credentials endpoint does
not permit reading a single entry.

Authentication

The ImprovMX REST API currently uses a simple authentication scheme. The
package will take care of setting this authentication scheme for users. Simply
construct a Session with the desired token, and the API will (assuming a valid
API Token) "just work"

Pagination

Several ImprovMX REST API endpoints currently use pagination to receive
results. At this time, the golang API does not support preventing pagination,
nor does it support manual pagination.

This will be fixed in the future.

*/
package improvmx
