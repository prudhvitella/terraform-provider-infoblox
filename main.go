package main

import (
	"github.com/hashicorp/terraform/plugin"

	infoblox "github.com/defilan/terraform-provider-infoblox/infoblox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: infoblox.Provider,
	})
}
