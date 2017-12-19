---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dns_zone_v2"
sidebar_current: "docs-opentelekomcloud-resource-dns-zone-v2"
description: |-
  Manages a DNS zone in the OpenTelekomCloud DNS Service
---

# opentelekomcloud\_dns\_zone_v2

Manages a DNS zone in the OpenTelekomCloud DNS Service.

## Example Usage

### Automatically detect the correct network

```hcl
resource "opentelekomcloud_dns_zone_v2" "example.com" {
  name = "example.com."
  email = "jdoe@example.com"
  description = "An example zone"
  ttl = 3000
  zone_type = "public"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    Keypairs are associated with accounts, but a Compute client is needed to
    create one. If omitted, the `region` argument of the provider is used.
    Changing this creates a new DNS zone.

* `name` - (Required) The name of the zone. Note the `.` at the end of the name.
  Changing this creates a new DNS zone.

* `email` - (Optional) The email contact for the zone record.

* `zone_type` - (Optional) The type of zone. Can either be `public` or `private`.
  Changing this creates a new zone.

* `router` - (Optional) Router configuration block which is required if zone_type is private.

* `ttl` - (Optional) The time to live (TTL) of the zone.

* `description` - (Optional) A description of the zone.

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new zone.

The `router` block supports:

* `router_id` - (Required) The router UUID. Changing this creates a new zone.

* `router_region` - (Required) The region of the router. Changing this creates a new zone.



## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `email` - See Argument Reference above.
* `zone_type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `masters` - An array of master DNS servers.
* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the zone ID:

```
$ terraform import opentelekomcloud_dns_zone_v2.zone_1 <zone_id>
```
