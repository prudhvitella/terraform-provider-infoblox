provider "infoblox" {
    username = "${var.infoblox_username}"
    password = "${var.infoblox_password}"
    host  = "https://infoblox.alaskaair.com"
    sslverify = false
    usecookies = false
}

resource "infoblox_record_host" "host" {
    name = "seadvmaherinfotest007"
    comment = "Bozo test"
    ipv4addr {
    address = "10.80.102.48"
  }
    configure_for_dns = false
}

/* resource "infoblox_ip" "ip" { */
/*   cidr = "10.80.100.0/22" */
/* } */

/* resource "infoblox_record_host" "test" { */
/*   name = "test_name" */
/*   ipv4addr { */
/*     address = "${infoblox_ip.ip.ipaddress}" */
/*   } */
/*   configure_for_dns = false */
/*   comment = "test comment" */
/* } */

/* output "ip" { */
/*  value = "${infoblox_ip.ip.ipaddress}" */
/* } */
