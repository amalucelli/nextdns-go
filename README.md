# nextdns-go

Go client library for [NextDNS](https://nextdns.io/) API.

## Install

```bash
go get github.com/amalucelli/nextdns-go/nextdns
```

## Requirements

An API Key is required to interact with the NextDNS API.
You can find your API Key in the [NextDNS account](https://my.nextdns.io/account) page.

## API

The [official API documentation](https://nextdns.github.io/api/) was the base document for this package.

## Usage

Here is an example usage of the NextAPI Go client for the `/profiles` endpoint:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/amalucelli/nextdns-go/nextdns"
)

func main() {
	// get the api key from the environment
	key := os.Getenv("NEXTDNS_API_KEY")

	// client client with a custom API key
	ctx := context.Background()
	client, _ := nextdns.New(
		nextdns.WithAPIKey(key),
	)

	// create a new profile
	id, _ := client.Profiles.Create(ctx, &nextdns.CreateProfileRequest{})

	// set a few settings like the name and some other attributes
	patch := nextdns.Profile{
		Name: "nextdns-go-example",
		Settings: nextdns.ProfileSettings{
			Web3: true,
		},
	}

	// update the profile
	_ = client.Profiles.Patch(ctx, &nextdns.PatchProfileRequest{Profile: id}, patch)

	// get the profile details to check the settings
	profile, _ := client.Profiles.Get(ctx, &nextdns.GetProfileRequest{
		Profile: id,
	})
	fmt.Printf("%q profile name: %s\n", id, profile.Name)
	fmt.Printf("%q web3 status: %t\n", id, profile.Settings.Web3)

	// list all the profiles
	profiles, _ := client.Profiles.List(ctx, &nextdns.ListProfileRequest{})
	fmt.Printf("Found %d profiles\n", len(profiles))
	for _, p := range profiles {
		fmt.Printf("ID: %q\n", p.ID)
		fmt.Printf("Name: %q\n", p.Name)
	}

	// delete the profile
	_ = client.Profiles.Delete(ctx, &nextdns.DeleteProfileRequest{
		Profile: id,
	})
}
```
