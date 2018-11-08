# [Terraform](https://github.com/hashicorp/terraform) Infoblox Provider

[![Build
status](https://travis-ci.org/prudhvitella/terraform-provider-infoblox.svg)](https://travis-ci.org/prudhvitella/terraform-provider-infoblox)

The Infoblox provider is used to interact with the
resources supported by Infoblox. The provider needs to be configured
with the proper credentials before it can be used.

## Download

Download builds for Darwin, Linux and Windows from the [releases page](https://github.com/prudhvitella/terraform-provider-infoblox/releases/).

## Example Usage

```hcl
# Configure the Infoblox provider
provider "infoblox" {
    username = "${var.infoblox_username}"
    password = "${var.infoblox_password}"
    host  = "${var.infoblox_host}"
    sslverify = "${var.infoblox_sslverify}"
    usecookies = "${var.infoblox_usecookies}"
}

# Create a record
resource "infoblox_record_a" "www" {
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

# infoblox\_record\_host

Provides an Infoblox Host record resource.

## Example Usage

```hcl
resource "infoblox_record_host" "host" {
  name              = "terraformhost.platform.test-aib.pri"
  configure_for_dns = false

  ipv4addr {
    address = "10.89.130.30"
  }

  ipv4addr {
    address            = "10.89.130.31"
    configure_for_dhcp = true
    mac                = "01-23-45-67-89-10"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the record
* `ipv4addr` - (Required) An IPv4 address object. At least one `iv4addr` or `ipv6addr` must be specified. See [ipv4addr options](#Ipv4addr_options) below.
* `ipv6addr` - (Required) An IPv6 address object. At least one `iv4addr` or `ipv6addr` must be specified. See [ipv6addr options](#Ipv6addr_options) below.
* `configure_for_dns` - (Boolean, Optional) Specify whether DNS should be configured for the record; defaults to `false`
* `comment` - (Optional) The comment for the record
* `ttl` - (Integer, Optional) The TTL of the record
* `view` - (Optional) The view of the record

### Ipv4 options

* `address` - (Required) The IPv4 address of the object
* `configure_for_dhcp` - (Boolean, Optional) Specifies whether the IPv4 address object should be configured for DHCP
* `mac` - (Optional) The MAC address of the resource

### Ipv6 options

* `address` - (Required) The IPv6 address of the object
* `configure_for_dhcp` - (Boolean, Optional) Specifies whether the IPv4 address object should be configured for DHCP
* `mac` - (Optional) The MAC address of the resource

# infoblox\_record\_a

Provides an Infoblox A record resource.

## Example Usage

```hcl
resource "infoblox_record_a" "web" {
  address = "10.1.2.3"
  name    = "some.fqdn.lan"

  comment = "ipv4 address for Acme web server"
  ttl     = 3600
  view    = "default"
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Required) The IPv4 address of the record
* `name` - (Required) The FQDN of the record
* `comment` - (Optional) The comment for the record
* `ttl` - (Integer, Optional) The TTL of the record
* `view` - (Optional) The view of the record

# infoblox\_record\_aaaa

Provides an Infoblox AAAA record resource.

## Example Usage

```hcl
resource "infoblox_record_aaaa" "web" {
  address = "2001:db8:85a3::8a2e:370:7334"
  name    = "some.fqdn.lan"

  comment = "ipv6 address for Acme web server"
  ttl     = 3600
  view    = "default"
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Required) The IPv6 address of the record
* `name` - (Required) The FQDN of the record
* `comment` - (Optional) The comment for the record
* `ttl` - (Integer, Optional) The TTL of the record
* `view` - (Optional) The view of the record

# infoblox\_record\_cname

Provides an Infoblox CNAME record resource.

## Example Usage

```hcl
resource "infoblox_record_cname" "www" {
  canonical = "fqdn.lan"
  name      = "www.fqdn.lan"

  comment = "ipv6 address for Acme web server"
  ttl     = 3600
  view    = "www.fqdn.lan is an alias for fqdn.lan"
}
```

## Argument Reference

The following arguments are supported:

* `canonical` - (Required) The canonical address to point to
* `name` - (Required) The FQDN of the alias
* `comment` - (Optional) The comment for the record
* `ttl` - (Integer, Optional) The TTL of the record
* `view` - (Optional) The view of the record

# infoblox\_record\_ptr

Provides an Infoblox PTR record resource.

## Example Usage

```hcl
resource "infoblox_record_ptr" "ptr" {
  ptrdname = "some.fqdn.lan"
  address  = "10.0.0.10.in-addr.arpa"

  comment = "Reverse lookup for some.fqdn.lan"
  ttl     = 3600
  view    = "default"
}
```

## Argument Reference

The following arguments are supported:

* `ptrdname` - (Required) The
* `address` - (Required, conflicts with `name`) This field is required if you do not use the name field. Either the IP address or name is required. Example: 10.0.0.11. If the PTR record belongs to a forward-mapping zone, this field is empty. Accepts both IPv4 and IPv6 addresses.
* `name` - (Required, conflicts with `address`) This field is required if you do not use the address field. Either the IP address or name is required. Example: 10.0.0.10.in.addr.arpa
* `comment` - (Optional) The comment for the record
* `ttl` - (Integer, Optional) The TTL of the record
* `view` - (Optional) The view of the record

# infoblox\_record\_txt

Provides an Infoblox TXT record resource.

## Example Usage

```hcl
resource "infoblox_record_txt" "txt" {
  name = "some.fqdn.lan"
  text = "Welcome to the Jungle"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)  The name of the TXT record
* `text` - (Required) The text of the TXT record
* `comment` - (Optional) The comment for the record
* `ttl` - (Integer, Optional) The TTL of the record
* `view` - (Optional) The view of the record

# infoblox\_record\_srv

Provides an Infoblox SRV record resource.

## Example Usage

```hcl
resource "infoblox_record_srv" "srv" {
  name = "bind_srv.domain.com"
  port = 1234
  priority = 1
  weight = 1
  target = "old.target.test.org"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)  The name of the record
* `port` - (Integer, Required) The port of the SRV record
* `priority` - (Integer, Required) The priority of the SRV record
* `weight` - (Integer, Required) The weight of the SRV record
* `target` - (Required) The target of the SRV record
* `comment` - (Optional) The comment for the record
* `ttl` - (Integer, Optional) The TTL of the record
* `view` - (Optional) The view of the record

# infoblox\_ip

Queries the next available IP address from a network and returns it in a computed variable
that can be used by the infoblox_record resource.

## Example Usage

```hcl
# Acquire the next available IP from a network CIDR
# it will create a variable called "ipaddress"
resource "infoblox_ip" "ip" {
  cidr = "10.0.0.0/24"
}

resource "infoblox_record_a" "web" {
  address = "${infoblox_ip.ip.ipaddress}"
  name    = "some.fqdn.lan"

  comment = "ipv4 address for Acme web server"
  ttl     = 3600
  view    = "default"
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

# Acquire free IP address from within a specific
# range of addresses
resource "infoblox_ip" "ipAddressFromRange" {
  ip_range = "10.0.0.20-10.0.0.60"
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) The network to search for - example 10.0.0.0/24. Cannot be specified with `ip_range`
* `exclude` - (Optional) A list of IP addresses to exclude
* `ip_range` - (Required) The IP range to search within - example 10.0.0.20-10.0.0.40. Cannot be
  specified with `cidr`

# Deprecated Resources

The following resources are deprecated and will no longer see active development. It is recommended you use the dedicated `infoblox_record_*` resources instead.

## infoblox\_record

Provides a Infoblox record resource.

### Example Usage

```hcl
# Add a record to the domain
resource "infoblox_record" "foobar" {
  value = "192.168.0.10"
  name = "terraform"
  domain = "mydomain.com"
  type = "A"
  ttl = 3600
}
```

### Argument Reference

See [related part of Infoblox Docs](https://godoc.org/github.com/fanatic/go-infoblox) for details about valid values.

The following arguments are supported:

* `domain` - (Required) The domain to add the record to
* `value` - (Required) The value of the record; its usage will depend on the `type` (see below)
* `name` - (Required) The name of the record
* `ttl` - (Integer, Optional) The TTL of the record
* `type` - (Required) The type of the record
* `comment` - (Optional) The comment of the record

### DNS Record Types

The type of record being created affects the interpretation of the `value` argument.

#### A Record

* `value` is the IPv4 address

#### CNAME Record

* `value` is the alias name

#### AAAA Record

* `value` is the IPv6 address

### Attributes Reference

The following attributes are exported:

* `domain` - The domain of the record
* `value` - The value of the record
* `name` - The name of the record
* `type` - The type of the record
* `ttl` - The TTL of the record

