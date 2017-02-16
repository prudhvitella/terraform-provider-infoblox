package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/prudhvitella/terraform-provider-infoblox/infoblox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: infoblox.Provider,
	})
}
