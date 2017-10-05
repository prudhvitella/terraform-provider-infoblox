[![Build
status](https://travis-ci.org/prudhvitella/terraform-provider-infoblox.svg)](https://travis-ci.org/prudhvitella/terraform-provider-infoblox)

# [Terraform](https://github.com/hashicorp/terraform) Infoblox Provider

The Infoblox provider is used to interact with the
resources supported by Infoblox. The provider needs to be configured
with the proper credentials before it can be used.

##  Download
Download builds for Darwin, Linux and Windows from the [releases page](https://github.com/prudhvitella/terraform-provider-infoblox/releases/).

## Example Usage

```
# Configure the Infoblox provider
provider "infoblox" {
    username = "${var.infoblox_username}"
    password = "${var.infoblox_password}"
    host  = "${var.infoblox_host}"
    sslverify = "${var.infoblox_sslverify}"
    usecookies = "${var.infoblox_usecookies}"
}

# Create a record
resource "infoblox_record" "www" {
    ...
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The Infoblox username. It must be provided, but it can also be sourced from the `INFOBLOX_USERNAME` environment variable.
* `password` - (Required) The password associated with the username. It must be provided, but it can also be sourced from the `INFOBLOX_PASSWORD` environment variable.
* `host` - (Required) The base url for the Infoblox REST API, but it can also be sourced from the `INFOBLOX_HOST` environment variable.
* `sslverify` - (Required) Enable ssl for the REST api, but it can also be sourced from the `INFOBLOX_SSLVERIFY` environment variable.
* `usecookies` - (Optional) Use cookies to connect to the REST API, but it can also be sourced from the `INFOBLOX_USECOOKIES` environment variable

# infoblox\_record

Provides a Infoblox record resource.

## Example Usage

```
# Add a record to the domain
resource "infoblox_record" "foobar" {
	value = "192.168.0.10"
	name = "terraform"
	domain = "mydomain.com"
	type = "A"
	ttl = 3600
}
```

## Argument Reference

See [related part of Infoblox Docs](https://godoc.org/github.com/fanatic/go-infoblox) for details about valid values.

The following arguments are supported:

* `domain` - (Required) The domain to add the record to
* `value` - (Required) The value of the record; its usage will depend on the `type` (see below)
* `name` - (Required) The name of the record
* `ttl` - (Integer, Optional) The TTL of the record
* `type` - (Required) The type of the record

## DNS Record Types

The type of record being created affects the interpretation of the `value` argument.

#### A Record

* `value` is the IPv4 address

#### CNAME Record

* `value` is the alias name

#### AAAA Record

* `value` is the IPv6 address

## Attributes Reference

The following attributes are exported:

* `domain` - The domain of the record
* `value` - The value of the record
* `name` - The name of the record
* `type` - The type of the record
* `ttl` - The TTL of the record

# infoblox\_ip

Queries the next available IP address from a network and returns it in a computed variable
that can be used by the infoblox_record resource.

## Example Usage

```
# Acquire the next available IP from a network CIDR
# it will create a variable called "ipaddress"
resource "infoblox_ip" "theIPAddress" {
	cidr = "10.0.0.0/24"
}


# Add a record to the domain
resource "infoblox_record" "foobar" {
	value = "${infoblox_ip.theIPAddress.ipaddress}"
	name = "terraform"
	domain = "mydomain.com"
	type = "A"
	ttl = 3600
}

# Exclude specific IP addresses when acquiring next
# avaiable IP from a network CIDR
resource "infoblox_ip" "excludedIPAddress" {
    cidr = "10.0.0.0/24"

    exclude = [
        "10.0.0.1",
        "10.0.0.2"
        # etc.
    ]
}

# Acquire gree IP address from within a specific
# range of addresses
resource "infoblox_ip" "ipAddressFromRange" {
    ip_range = "10.0.0.20-10.0.0.60"
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) The network to search for - example 10.0.0.0/24. Cannot be specified with `ip\_range`
* `exclude` - (Optional) A list of IP addresses to exclude
* `ip_range` - (Required) The IP range to search within - example 10.0.0.20-10.0.0.40. Cannot be
  specified with `cidr`
