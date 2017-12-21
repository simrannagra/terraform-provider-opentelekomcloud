resource "opentelekomcloud_dns_zone_v2" "dnszone" {
  count = "${var.dnszone != "" ? 1 : 0}"
  name  = "${var.dnszone}."
  email = "email@${var.dnszone}"
  ttl   = 6000
  zone_type = "public"
}

resource "opentelekomcloud_dns_recordset_v2" "recordset" {
  count   = "${var.dnszone != "" ? 1 : 0}"
  zone_id = "${opentelekomcloud_dns_zone_v2.dnszone.id}"
  name    = "${var.dnsname}.${var.dnszone}."
  ttl     = 3000
  type    = "A"
  records = ["${opentelekomcloud_networking_floatingip_v2.fip.*.address}"]
}

resource "opentelekomcloud_dns_zone_v2" "dnszone_private" {
  count = "${var.dnszone != "" ? 1 : 0}"
  name  = "local.${var.dnszone}."
  email = "email@${var.dnszone}"
  ttl   = 6000
  zone_type = "private"
  router {
    router_id = "${opentelekomcloud_networking_router_v2.router.id}"
    router_region = "eu-de"
  }
}
