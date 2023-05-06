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

APIs supported by this package:

- [x] Profile (`/profiles` and `/profiles/:profile`)
- [x] Analytics (`/profiles/:profile/analytics`)
- [ ] Logs (`/profiles/:profile/logs`)

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
	create := &nextdns.CreateProfileRequest{
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
	id, _ := client.Profiles.Create(ctx, create)

	// set a few settings like the name and some other attributes
	update := &nextdns.UpdateProfileRequest{
		ProfileID: id,
		Profile: &nextdns.Profile{
			Name: "nextdns-go-updated",
			Settings: &nextdns.Settings{
				Logs: &nextdns.SettingsLogs{
					Enabled: false,
				},
			},
		},
	}

	// update the profile
	_ = client.Profiles.Update(ctx, update)

	// get the profile details to check the settings
	profile, _ := client.Profiles.Get(ctx, &nextdns.GetProfileRequest{
		ProfileID: id,
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
		ProfileID: id,
	})
}
```

It's also possible to update directly the API child endpoints, like the `/profiles/:profile/denylist` endpoint:

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

	// set the profile id
	id := "abc123"

	// set the request to update the denylist
	request := &nextdns.UpdateDenylistRequest{
		ProfileID: id,
		ID:      "google.com",
		Denylist: &nextdns.Denylist{
			Active: true,
		},
	}

	// update the denylist
	_ = client.Denylist.Update(ctx, request)

	// list all the denylist entries
	list, _ := client.Denylist.Get(ctx, &nextdns.GetDenylistRequest{ProfileID: id})
	for _, p := range list {
		fmt.Printf("ID: %q\n", p.ID)
		fmt.Printf("Status: %t\n", p.Active)
	}
}
```
