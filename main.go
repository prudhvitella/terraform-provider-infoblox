package main

import (
	"github.com/hashicorp/terraform/plugin"
	"./infoblox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: infoblox.Provider,
	})
}
