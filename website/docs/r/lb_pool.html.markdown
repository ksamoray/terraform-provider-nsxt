---
layout: "nsxt"
page_title: "NSXT: nsxt_lb_pool"
sidebar_current: "docs-nsxt-resource-lb-pool"
description: |-
  Provides a resource to configure lb pool on NSX-T manager
---

# nsxt_lb_pool

Provides a resource to configure lb pool on NSX-T manager

## Example Usage

```hcl
resource "nsxt_lb_icmp_monitor" "lb_icmp_monitor" {
  display_name = "lb_icmp_monitor"
  fall_count   = 3
  interval     = 5
}

resource "nsxt_lb_pool" "lb_pool" {
  description              = "lb_pool provisioned by Terraform"
  display_name             = "lb_pool"
  algorithm                = "WEIGHTED_ROUND_ROBIN"
  min_active_members       = 1
  tcp_multiplexing_enabled = false
  tcp_multiplexing_number  = 3
  active_monitor_id        = "${nsxt_lb_icmp_monitor.lb_icmp_monitor.id}"
  snat_translation_type    = "LbSnatAutoMap"
 
  member {
    admin_state                = "ENABLED"
    backup_member              = "false"
    display_name               = "1st-member"
    ip_address                 = "1.1.1.1"
    max_concurrent_connections = "1"
    port                       = "87"
    weight                     = "1"
  }

  tag = {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `active_monitor_id` - (Optional) Active health monitor Id. If one is not set, the active healthchecks will be disabled.
* `algorithm` - (Optional) Load balancing algorithm controls how the incoming connections are distributed among the members. Supported algorithms are: ROUND_ROBIN, WEIGHTED_ROUND_ROBIN, LEAST_CONNECTION, WEIGHTED_LEAST_CONNECTION, IP_HASH.
* `member` - (Optional) Server pool consists of one or more pool members. Each pool member is identified, typically, by an IP address and a port. Each member has the following arguments:
  * `admin_state` - (Optional) Pool member admin state.
  * `backup_member` - (Optional) A boolean flag which reflects whether this is a backup pool member. Backup servers are typically configured with a sorry page indicating to the user that the application is currently unavailable. While the pool is active (a specified minimum number of pool members are active) BACKUP members are skipped during server selection. When the pool is inactive, incoming connections are sent to only the BACKUP member(s).
  * `display_name` - (Optional) The display name of this resource. pool member name.
  * `ip_address` - (Required) Pool member IP address.
  * `max_concurrent_connections` - (Optional) To ensure members are not overloaded, connections to a member can be capped by the load balancer. When a member reaches this limit, it is skipped during server selection. If it is not specified, it means that connections are unlimited.
  * `port` - (Optional) If port is specified, all connections will be sent to this port. Only single port is supported. If unset, the same port the client connected to will be used, it could be overrode by default_pool_member_port setting in virtual server. The port should not specified for port range case.
  * `weight` - (Optional) Pool member weight is used for WEIGHTED_ROUND_ROBIN balancing algorithm. The weight value would be ignored in other algorithms.
* `min_active_members` - (Optional) The minimum number of members for the pool to be considered active. This value is 1 by default.
* `passive_monitor_id` - (Optional) Passive health monitor Id. If one is not set, the passive healthchecks will be disabled.
* `snat_translation_type` - (Optional) Type of SNAT performed to ensure reverse traffic from the server can be received and processed by the loadbalancer. Supported types are: LbSnatAutoMap, Transparent
* `tcp_multiplexing_enabled` - (Optional) TCP multiplexing allows the same TCP connection between load balancer and the backend server to be used for sending multiple client requests from different client TCP connections. Disabled by default.
* `tcp_multiplexing_number` - (Optional) The maximum number of TCP connections per pool that are idly kept alive for sending future client requests. The default value for this is 6.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb pool.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb pool.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb pool can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_lb_pool.lb_pool UUID
```

The above would import the lb pool named `lb_pool` with the nsx id `UUID`