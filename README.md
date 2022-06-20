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

	// set a few settings like the name and some other attributes
	new := &nextdns.CreateProfileRequest{
		Name: "nextdns-go",
		Denylist: []*nextdns.Denylist{
			{
				ID:     "google.com",
				Active: true,
			},
			{
				ID:     "bing.com",
				Active: true,
			},
		},
		Allowlist: []*nextdns.Allowlist{
			{
				ID:     "duckduckgo.com",
				Active: true,
			},
			{
				ID:     "search.brave.com",
				Active: false,
			},
		},
		ParentalControl: &nextdns.ParentalControl{
			Categories: []*nextdns.ParentalControlCategories{
				{
					ID:     "gambling",
					Active: true,
				},
			},
		},
		Security: &nextdns.Security{
			AiThreatDetection: true,
		},
		Settings: &nextdns.Settings{
			Logs: &nextdns.SettingsLogs{
				Enabled: true,
			},
			Web3: true,
		},
	}

	// create a new profile
	id, _ := client.Profiles.Create(ctx, new)

	// set a few settings like the name and some other attributes
	settings := &nextdns.Profile{
		Name: "nextdns-go-updated",
		Settings: &nextdns.Settings{
			Logs: &nextdns.SettingsLogs{
				Enabled: false,
			},
		},
	}

	// update the profile
	_ = client.Profiles.Update(ctx, &nextdns.UpdateProfileRequest{Profile: id}, settings)

	// get the profile details to check the settings
	profile, _ := client.Profiles.Get(ctx, &nextdns.GetProfileRequest{
		Profile: id,
	})
	fmt.Printf("%q profile name: %s\n", id, profile.Name)
	fmt.Printf("%q logs status: %t\n", id, profile.Settings.Logs.Enabled)

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
